package mothership

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/mnordsletten/lotto/util"
	"github.com/sirupsen/logrus"
)

type InstanceHealth struct {
	Online       bool   `json:"Online"`
	TotalPanics  int    `json:"Panics"`
	IosVersion   string `json:"Version"`
	Time         string
	NewPanics    int
	PanicContent string
}

func (m *Mothership) CheckInstanceHealth() InstanceHealth {
	m.alias = "a"
	i, err := m.getInstanceInfo(m.alias)
	if err != nil {
		logrus.Warningf("Could not get instance info: %v", err)
	}

	// Check for any crashes, and get the latest crash output if necessary
	crashIDs, err := m.getAllCrashesArray()
	if err != nil {
		logrus.Warningf("could not get crashes array: %v", err)
		return i
	}
	logrus.Debugf("All crashes: %v", crashIDs)
	i.NewPanics = len(crashIDs)
	if len(crashIDs) > 0 {
		i.PanicContent, err = m.getLatestCrashOutput(crashIDs)
		if err != nil {
			logrus.Warningf("Could not get panic output: %v", err)
		}
	}
	return i
}

func ConvertHealthToPrintableOutput(iHealth InstanceHealth, filename string) error {
	s := reflect.ValueOf(&iHealth).Elem()
	typeOfI := s.Type()

	x := make([][]string, 2)
	for i := 0; i < s.NumField(); i++ {
		// Headers first
		x[0] = append(x[0], typeOfI.Field(i).Name)

		// Then content
		f := s.Field(i)
		content := fmt.Sprintf("%v", f.Interface())
		x[1] = append(x[1], content)
	}
	if err := util.OutputWriter(x, filename); err != nil {
		logrus.Warningf("could not write instance health: %v", err)
	}

	return nil
}

func (m *Mothership) getInstanceInfo(id string) (InstanceHealth, error) {
	var i InstanceHealth
	i.Time = time.Now().Format(time.RFC3339)

	request := fmt.Sprintf("inspect-instance %s -o json", id)
	output, err := m.bin(request)
	if err != nil {
		return i, err
	}
	if err := json.Unmarshal([]byte(output), &i); err != nil {
		return i, fmt.Errorf("error unmarshaling instanceInfo: %v", err)
	}

	return i, nil
}

func (m *Mothership) getLatestCrashOutput(crashIDs []string) (string, error) {
	logrus.Debugf("Getting latest panic output")

	latestCrash := crashIDs[len(crashIDs)-1]
	logrus.Debugf("latest crashID: %s", latestCrash)
	crashContent, err := m.getSingleCrashOutput(latestCrash)
	if err != nil {
		logrus.Warningf("could not get crash output for %s: %v", latestCrash, err)
	}

	m.deleteCrashes(crashIDs)
	return crashContent, nil
}

func (m *Mothership) getAllCrashesArray() ([]string, error) {
	type crashes []string
	var c crashes
	request := fmt.Sprintf("instance-crashes %s -o json", m.alias)
	output, err := m.bin(request)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(output), &c); err != nil {
		return nil, fmt.Errorf("error unmarshaling instance-crashes: %v", err)
	}
	return c, nil
}

func (m *Mothership) getSingleCrashOutput(crashID string) (string, error) {
	request := fmt.Sprintf("instance-crash %s %s", m.alias, crashID)
	output, err := m.bin(request)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (m *Mothership) deleteCrashes(crashNames []string) {
	for _, crashID := range crashNames {
		request := fmt.Sprintf("delete-crash %s %s", m.alias, crashID)
		_, err := m.bin(request)
		if err != nil {
			logrus.Warningf("could not delete crash %s: %v", crashID, err)
		}
	}
}