// Package config provides application support for YAML config files.
package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration stores the configuration for this app.
type Configuration struct {
	Instance    GitlabInstance     `json:"gitlab-instance"`
	QueryConfig QueryConfiguration `json:"query-configuration"`
}

// GitlabInstance is the instance of Gitlab that will be accessed.
type GitlabInstance struct {
	URL         string `json:"url"`
	Username    string `json:"username"`
	AccessToken string `json:"access-token"`
	Project     string `json:"project"`
}

// QueryConfiguration allows the gitlab note queries to be configured.
type QueryConfiguration struct {
	TimeEntryDetectionRegex  string `json:"time-detection-regex"`
	TimeEntryExtractionRegex string `json:"time-extraction-regex"`
}

// Default will generate the default configuration.
func Default() Configuration {
	return Configuration{
		Instance: GitlabInstance{
			URL:         "gitlab.somedomain.com",
			Username:    "coder1",
			AccessToken: "12345",
			Project:     "collection/project",
		},
		QueryConfig: QueryConfiguration{
			TimeEntryDetectionRegex: `/time spent/g`,
			TimeEntryExtractionRegex: `(?P<month>[0-9]+)(mo)|(?P<week>[0-9]+)(w)|(?P<day>[0-9]+)(d)|(?P<hour>[0-9]+)(h)|(?P<minute>[0-9]+)(m)|(?P<second>[0-9]+)(s)`,
		},
	}
}

// Load will retrieve a configuration from a yaml file.
func Load(path string) (*Configuration, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data Configuration
	err = yaml.Unmarshal(fileBytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// Save will write a configuration to yaml file.
func (c *Configuration) Save(path string) error {
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, y, 0644)
	if err != nil {
		return err
	}
	return nil
}
