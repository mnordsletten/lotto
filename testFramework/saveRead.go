package testFramework

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func (t *TestConfig) SaveToDisk() error {
	// Create data dir if it does not exist
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	// Save instance to file
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	filename := path.Join("data", t.ID+".json")
	if err := ioutil.WriteFile(filename, b, 0600); err != nil {
		return err
	}
	return nil
}

// ReadFromDisk takes a path and reads the testspec and returns a TestConfig
func ReadFromDisk(servicePath string) (*Service, error) {
	test := &Service{}
	test.ServicePath = servicePath
	test.Name = path.Base(test.ServicePath)
	/*
		if err := verifyTestFiles(testPath); err != nil {
			return test, err
		}
	*/
	testSpec, err := ioutil.ReadFile(path.Join(test.ServicePath, "testspec.json"))
	if err != nil {
		return test, fmt.Errorf("error reading test file %s: %v", test.ServicePath, err)
	}
	if err = json.Unmarshal(testSpec, test); err != nil {
		return test, fmt.Errorf("error decoding json: %v", err)
	}

	// append testPath to all files
	if test.NaclFile != "" {
		test.NaclFile = path.Join(test.ServicePath, test.NaclFile)
	}
	if test.CustomServicePath != "" {
		dir, err := os.Executable()
		if err != nil {
			return test, fmt.Errorf("error getting current dir: %v", err)
		}
		// Requires full path due to getting mounted into docker
		test.CustomServicePath = path.Join(path.Dir(dir), test.ServicePath, test.CustomServicePath)
	}
	return test, nil
}

// verifyTestFiles checks for the existence of selected files in the path supplied
func verifyTestFiles(testPath string) error {
	expectedFiles := []string{"testspec.json", "*.sh"}

	for _, file := range expectedFiles {
		pathToCheck := path.Join(testPath, file)
		matches, err := filepath.Glob(pathToCheck)
		if err != nil {
			return fmt.Errorf("error looking for file: %s, %v", file, err)
		}
		if len(matches) < 1 {
			return fmt.Errorf("%s not found", file)
		}
	}

	return nil
}
