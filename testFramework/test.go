package testFramework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"text/template"
	"time"

	"github.com/mnordsletten/lotto/environment"
	"github.com/sirupsen/logrus"
)

type Test struct {
	Name                string
	TestPath            string
	ClientCommandScript string                 `json:"clientcommandscript"`
	Setup               environment.SSHClients `json:"setup"`
	Cleanup             environment.SSHClients `json:"cleanup"`
}

func (t *Test) Run(env environment.Environment, templates []TemplateKeyValue) (TestResult, error) {
	var result TestResult
	result.Name = t.Name
	// Prepare clients/servers for test
	defer func() {
		if err := t.cleanupTestMachines(env); err != nil {
			logrus.Warningf("error cleaning up after test: %v", err)
		}
	}()
	if err := t.prepareTestMachines(env); err != nil {
		return result, fmt.Errorf("error preparing test machines: %v", err)
	}
	f, err := ioutil.ReadFile(path.Join(t.TestPath, t.ClientCommandScript))
	if err != nil {
		return result, fmt.Errorf("error reading clientCommandScript: %v", err)
	}

	// Parse requires a string
	m, err := template.New("test").Parse(string(f))
	if err != nil {
		return result, fmt.Errorf("error parsing template: %v", err)
	}

	tem := struct {
		Template map[string]string
	}{}

	tem.Template = make(map[string]string, len(templates))

	for _, t := range templates {
		tem.Template[t.Key] = t.Value
	}

	var script bytes.Buffer
	if err = m.Execute(&script, tem); err != nil {
		return result, fmt.Errorf("error executing template: %v", err)
	}
	start := time.Now()
	testOutput, err := env.RunClientCmdScript(1, script.Bytes())
	if err != nil {
		return result, fmt.Errorf("error running client command script: %v", err)
	}
	result.Duration = time.Now().Sub(start)

	// Parse test results
	if err = json.Unmarshal(testOutput, &result); err != nil {
		return result, fmt.Errorf("could not parse testResults: %v\ntestOutput: %s", err, testOutput)
	}
	result.SuccessPercentage = float32(result.Received) / float32(result.Sent) * 100

	return result, nil
}

func (t *Test) prepareTestMachines(env environment.Environment) error {
	logrus.Debug("Preparing test machines")
	if err := runScriptsOnClients(t.TestPath, env, t.Setup); err != nil {
		return fmt.Errorf("error preparing test: %v", err)
	}
	return nil
}

func (t *Test) cleanupTestMachines(env environment.Environment) error {
	logrus.Debug("Cleaning test machines")
	if err := runScriptsOnClients(t.TestPath, env, t.Cleanup); err != nil {
		return fmt.Errorf("error cleaning test: %v", err)
	}
	return nil
}

func runScriptsOnClients(scriptsPath string, env environment.Environment, scripts environment.SSHClients) error {
	for i := 1; i <= 4; i++ {
		scriptName, err := scripts.GetClientByInt(i)
		if err != nil {
			return fmt.Errorf("error getting client%d: %v", i, err)
		} else if scriptName == "" {
			continue
		}
		logrus.Debugf("running script: %s on client%d", scriptName, i)
		scriptPath := path.Join(scriptsPath, scriptName)
		fileByte, err := ioutil.ReadFile(scriptPath)
		if err != nil {
			return fmt.Errorf("error reading script file: %s: %v", scriptPath, err)
		}
		if output, err := env.RunClientCmdScript(i, fileByte); err != nil {
			return fmt.Errorf("error running script %s on client%d: out: %s: %v", scriptPath, i, output, err)
		}
	}
	return nil
}

func debugMode() {
	logrus.Info("Debug Mode. Waiting for ctrl-c to shut down")
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	logrus.Info("Debug Mode done")
}
