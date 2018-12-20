package mothership

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/mnordsletten/lotto/environment"
	"github.com/mnordsletten/lotto/util"
	"github.com/sirupsen/logrus"
)

// DeployNacl takes a nacl file name and deploys it to the Mothership
// Returns the naclId and the image checksum
func (m *Mothership) DeployNacl(naclFileName string) (string, string, error) {
	naclID, err := m.pushNacl(naclFileName)
	if err != nil {
		return naclID, "", err
	}

	checksum, err := m.build(naclID)
	if err != nil {
		return naclID, checksum, err
	}
	if err := m.deploy(checksum); err != nil {
		return naclID, checksum, err
	}
	return naclID, checksum, nil
}

// createCleanStarbase uses a standard nacl to build and pull down the image
func (m *Mothership) createCleanStarbase() error {
	naclFileName := "clean-starbase.nacl"
	var naclID, checksum string
	var err error
	// Push-NaCl
	if naclID, err = m.pushNacl(naclFileName); err != nil {
		return fmt.Errorf("error pushing NaCl: %v", err)
	}

	// Build
	if checksum, err = m.build(naclID); err != nil {
		return fmt.Errorf("error building %s: %v", naclFileName, err)
	}

	// Remove resources from the mothership
	defer func() {
		if err := m.DeleteNacl(naclID); err != nil {
			logrus.Errorf("error deleting clean-starbase nacl: %v", err)
		}
		if err := m.DeleteImage(checksum); err != nil {
			logrus.Errorf("error deleting clean-starbase image: %v", err)
		}
	}()

	// Pull image
	if err := m.pullImage(checksum, cleanStarbaseImage); err != nil {
		return fmt.Errorf("error pulling built image: %v", err)
	}

	return nil
}

// LaunchCleanStarbase launches a starbase in the specified environment and connects it
// to the Mothership.
func (m *Mothership) LaunchCleanStarbase(env environment.Environment) error {
	// Remove old alias, if it does not exist an error is returned but that does not matter
	if err := m.deleteInstanceByAlias(m.Alias); err != nil {
		logrus.Debugf("%s could not be removed: %v", m.Alias, err)
	}

	// Create a new tag which is inserted in the image that is built
	if err := m.createCleanStarbase(); err != nil {
		return err
	}

	// Launch the newly built starbase
	if err := m.Launch(cleanStarbaseImage, env); err != nil {
		return err
	}

	// Wait until the starbase with the "unique" tag we just created connects
	id, err := m.waitUntilStarbaseConnects(m.lastBuildTag)
	if err != nil {
		return err
	}

	// Set the alias to something generic
	if err := m.setAlias(m.Alias, id); err != nil {
		return err
	}
	return nil
}

func (m *Mothership) CheckStarbaseIDInUse() bool {
	request := fmt.Sprintf("inspect-instance %s -o json", m.Alias)
	response, err := m.bin(request)
	if err != nil {
		logrus.Infof("Mothership: starbase with alias: %s does not exist", m.Alias)
		return false
	}
	logrus.Infof("Mothership: starbase with alias: %s exists", m.Alias)
	var star starbase
	if err := json.Unmarshal([]byte(response), &star); err != nil {
		logrus.Debugf("could not check if starbase is connected: %v", err)
		return false
	}
	if star.Status != "connected" {
		logrus.Infof("Mothership: starbase with alias: %s is disconnected", m.Alias)
		return false
	}
	return true
}

func (m *Mothership) waitUntilStarbaseConnects(tag string) (string, error) {
	logrus.Debugf("Now waiting for starbase with tag: %s to connect", tag)
	var ids []struct {
		ID string `json:"id"`
	}
	deadLine := time.Now().Add(120 * time.Second)
	for time.Now().Before(deadLine) {
		req := fmt.Sprintf("search %s --instancefilter tag -o json", tag)
		if output, err := m.bin(req); err == nil {
			// Cmd did not return an error
			err := json.Unmarshal([]byte(output), &ids)
			if err != nil {
				return "", fmt.Errorf("error unmarshaling json from search: %v", err)
			}
			for _, id := range ids {
				// there is at least 1 key in the output
				logrus.Debugf("Starbase with tag: %s connected", tag)
				return id.ID, nil
			}
		}
		time.Sleep(2 * time.Second)
	}
	return "", fmt.Errorf("Starbase %s never connected to mothership", tag)
}

// BuildPushAndDeployCustomService will build a custom service using docker, push it and deploy if the argument
// deploy is not set to false
func (m *Mothership) BuildPushAndDeployCustomService(customServicePath, builderName string, noDeploy bool) (string, error) {
	imageName := path.Base(customServicePath)
	logrus.Infof("Building custom service %s", imageName)
	imagePath := path.Join(customServicePath, "build", imageName)
	if err := util.BuildServiceInDocker(customServicePath, m.uplinkFileName, builderName); err != nil {
		return "", fmt.Errorf("could not build custom service: %v", err)
	}

	imageID, err := m.PushImage(imagePath)
	if err != nil {
		return "", fmt.Errorf("could not push image %s: %v", imagePath, err)
	}
	if !noDeploy {
		if err := m.deploy(imageID); err != nil {
			return "", fmt.Errorf("error deploying custom service: %v", err)
		}
	}
	return imageID, nil
}

func (m *Mothership) PrepareBuilder(builderName string) error {
	logrus.Debugf("Preparing builder with name: %s", builderName)
	// First update all providers
	if err := m.BobProvidersUpdate(); err != nil {
		return fmt.Errorf("error updating providers: %v", err)
	}
	// Get a list of all available bobs
	logrus.Debugf("Getting list of all available bobs")
	bobsJson, err := m.BobsList()
	if err != nil {
		return fmt.Errorf("error getting list of all bobs")
	}
	// Find the ID and provider of the builder that we want
	var bobs []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Provider string `json:"providerId"`
	}
	if err = json.Unmarshal([]byte(bobsJson), &bobs); err != nil {
		return fmt.Errorf("error unmarshaling bobs: %v", err)
	}
	// Loop through all bobs looking for builderName
	logrus.Debugf("Looking for bob that matches builderName: %s", builderName)
	var builderID, providerID string
	var matchFound bool
	for _, bob := range bobs {
		if bob.Name == builderName {
			logrus.Debugf("Found bob that matches, builderID: %s, providerID: %s", bob.ID, bob.Provider)
			builderID = bob.ID
			providerID = bob.Provider
			matchFound = true
			break
		}
	}
	if !matchFound || len(builderID) == 0 || len(providerID) == 0 {
		return fmt.Errorf("problem with match: matchFound: %t builderID: %s providerID: %s", matchFound, builderID, providerID)
	}
	// Prepare the builder that was just found
	if err = m.BobPrepare(builderID, providerID); err != nil {
		return fmt.Errorf("error preparing bob with ID: %s and providerID: %s: %v", builderID, providerID, err)
	}
	logrus.Debugf("setting Mothership builderID to: %s", builderID)
	m.BuilderID = builderID
	return nil
}
