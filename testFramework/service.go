package testFramework

import (
	"fmt"

	"github.com/mnordsletten/lotto/environment"
	"github.com/mnordsletten/lotto/mothership"
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
	Path     string           `json:"path"`
	Template TemplateKeyValue `json:"template"`
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

func (s *Service) RunTest(iterations int, env environment.Environment, mother *mothership.Mothership) (TestResult, error) {
	fmt.Printf("Service: %+v\n", s)

	for _, test := range s.Tests {
		fmt.Printf("test: %+v\n", test)
	}
	return TestResult{}, nil
}
