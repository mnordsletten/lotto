package cmd

import (
	"fmt"

	"github.com/mnordsletten/lotto/environment"
	"github.com/mnordsletten/lotto/mothership"
	"github.com/mnordsletten/lotto/prettyoutput"
	"github.com/mnordsletten/lotto/testFramework"
)

func testProcedure(service *testFramework.Service, tests []string, env environment.Environment, mother *mothership.Mothership) (bool, error) {
	pretty := pretty.NewPrettyTest(service.Name)
	pretty.PrintHeader()
	pretty.PrintTable(service.StringSlice())

	// BUILD & DEPLOY. 3 options:
	// 1. Push NaCl and build on Mothership
	// 2. Build service locally using docker
	// 3. No building at all
	if err := build(service, mother); err != nil {
		return false, fmt.Errorf("error building: %v", err)
	}

	// TEST script. 2 options:
	// 1. HOSTcommandscript
	// 2. CLIENTcommandscript
	// Run client command
	// numRuns flag taken into account
	for _, testPath := range tests {
		result, err := service.RunTest(testPath, env, mother)
		if err != nil {
			return false, fmt.Errorf("error testing service %v", err)
		}
		// RESULTS print test results
		pretty.PrintTable(result.StringSlice())
		if !result.Success {
			fmt.Printf("Raw output: %s\n", result.Raw)
		}
	}

	/*
		// VERIFY starbase status
		health := mother.CheckInstanceHealth()
		logrus.Info(health)

		pretty.EndTest()
		reporting.SendReport(reporting.Dashboard{
			Address:           "http://localhost:7070/upload",
			MothershipVersion: "v1",
			IncludeOSVersion:  builderName,
			Environment:       cmdEnv,
			TestResult:        result,
		})
	*/
	//return result.Success, nil
	return true, nil
}

func getServicesToTest(possibleServices []string) ([]*testFramework.Service, error) {
	// Get the config for every service to be tested
	var services []*testFramework.Service
	for _, servicePath := range possibleServices {
		service, err := testFramework.ReadFromDisk(servicePath)
		if err != nil {
			return nil, fmt.Errorf("Could not read service spec: %v", err)
		}

		// enable debugMode if specified
		service.DebugMode = debugMode
		services = append(services, service)
	}
	return services, nil
}

func build(test *testFramework.Service, mother *mothership.Mothership) error {
	// Return early if skip rebuild has been set
	if test.SkipRebuild {
		return nil
	}
	var err error
	// Boot NaCl service to starbase, only if NaclFile is specified
	if test.NaclFile != "" {
		if test.NaclFileShasum, test.ImageID, err = mother.DeployNacl(test.NaclFile); err != nil {
			return fmt.Errorf("could not deploy: %v", err)
		}
	}
	// Build and deploy custom service if specified
	if test.CustomServicePath != "" {
		if test.ImageID, err = mother.BuildPushAndDeployCustomService(test.CustomServicePath, builderName, test.NoDeploy); err != nil {
			return fmt.Errorf("could not build and push custom service: %v", err)
		}
	}
	return nil
}
