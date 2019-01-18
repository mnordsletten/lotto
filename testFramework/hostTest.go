package testFramework

type HostCommandTemplate struct {
	MothershipBinPathAndName string
	OriginalAlias            string
	ImageID                  string
	BuilderID                string
}

/*
func (t *TestConfig) runHostTest(mother *mothership.Mothership) ([]byte, error) {
	// Process script file as template, replace template objects with actual info.
	logrus.Debugf("Running host test: %s", t.HostCommandScript)
	f, err := ioutil.ReadFile(t.HostCommandScript)
	if err != nil {
		return nil, fmt.Errorf("error reading lotto test template: %v", err)
	}

	// Parse requires a string
	m, err := template.New("test").Parse(string(f))
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	templ := HostCommandTemplate{
		MothershipBinPathAndName: mother.CLICommand(),
		OriginalAlias:            mother.Alias,
		ImageID:                  t.ImageID,
		BuilderID:                mother.BuilderID,
	}
	var script bytes.Buffer
	if err = m.Execute(&script, templ); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	out, err := util.ExternalCommandInput(html.UnescapeString(script.String()), nil)
	if err != nil {
		return out, fmt.Errorf("Host test external command failed: %v", err)
	}
	// Unmarshal test results into testResult
	return out, nil
}
*/
