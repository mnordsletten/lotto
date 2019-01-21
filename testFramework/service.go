package testFramework

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/mnordsletten/lotto/environment"
	"github.com/mnordsletten/lotto/mothership"
	"github.com/sirupsen/logrus"
)

type Service struct {
	NaclFile          string           `json:"naclfile"`
	Tests             []TestTypeConfig `json:"tests"`
	CustomServicePath string           `json:"customservicepath"`
	NoDeploy          bool             `json:"nodeploy"`
	NaclFileShasum    string
	SkipRebuild       bool
	Name              string
	ServicePath       string
	DebugMode         bool
	ImageID           string
}

type TestTypeConfig struct {
	Path     string             `json:"path"`
	Template []TemplateKeyValue `json:"template"`
}

type TemplateKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *Service) StringSlice() [][]string {
	var output [][]string
	output = append(output, []string{"Name", s.Name})
	if s.SkipRebuild {
		output = append(output, []string{"Skip rebuild", "[X]"})
	} else {
		output = append(output, []string{"Skip rebuild", "[ ]"})
	}
	return output
}

func (s *Service) RunTest(testPath string, env environment.Environment, mother *mothership.Mothership) (TestResult, error) {
	var testResult TestResult
	if ok, test := s.CheckIfTestIsConfigured(testPath); ok {
		t, err := readTestSpec(testPath)
		if err != nil {
			return testResult, fmt.Errorf("error reading test spec: %v", err)
		}
		testResult, err = t.Run(env, test.Template)
		if err != nil {
			return testResult, fmt.Errorf("error running test: %v", err)
		}
	}
	return testResult, nil
}

func (s *Service) CheckIfTestIsConfigured(testPath string) (bool, TestTypeConfig) {
	for _, test := range s.Tests {
		if strings.TrimSpace(testPath) == strings.TrimSpace(path.Join(s.ServicePath, test.Path)) {
			return true, test
		}
	}
	return false, TestTypeConfig{}
}

func readTestSpec(specPath string) (*Test, error) {
	t := &Test{}
	t.TestPath = specPath
	t.Name = filepath.Base(specPath)
	specFilePath := path.Join(specPath, "spec.json")
	testSpec, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return t, fmt.Errorf("error reading test %s: %v", specFilePath, err)
	}
	if err = json.Unmarshal(testSpec, t); err != nil {
		return t, fmt.Errorf("error decoding json: %v", err)
	}
	return t, nil
}

func (s *Service) cleanupService(mother *mothership.Mothership, env environment.Environment) {
	// Remove NaCl
	if len(s.NaclFileShasum) > 0 {
		if err := mother.DeleteNacl(s.NaclFileShasum); err != nil {
			logrus.Errorf("could not clean up nacl: %v", err)
		}
	}

	// Remove image
	if len(s.ImageID) > 0 {
		if err := mother.DeleteImage(s.ImageID); err != nil {
			logrus.Errorf("could not clean up image: %v", err)
		}
	}
}
