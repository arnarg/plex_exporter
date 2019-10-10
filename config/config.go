package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type PlexConfig struct {
	ListenAddress string             `yaml:"address" flag:"listen-address"`
	LogLevel      string             `yaml:"logLevel" flag:"log-level"`
	LogFormat     string             `yaml:"logFormat" flag:"format"`
	AutoDiscover  bool               `yaml:"autoDiscover" flag:"auto-discover"`
	Token         string             `yaml:"token" flag:"token"`
	Servers       []PlexServerConfig `yaml:"servers"`
}

type PlexServerConfig struct {
	BaseURL  string `yaml:"baseUrl"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`
}

func Load(c *cli.Context) (*PlexConfig, error) {
	plexConfig := &PlexConfig{}
	configPath := c.String("config-path")

	// Get absolute path of config file
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	// Check if the file already exists
	_, err = os.Stat(absPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		// Read config file
		yamlString, err := ioutil.ReadFile(absPath)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(yamlString, plexConfig)
		if err != nil {
			return nil, err
		}
	}

	// Merge config with cli flags
	MergeConfig(plexConfig, c)

	return plexConfig, nil
}

func MergeConfig(conf *PlexConfig, c *cli.Context) {
	// Overwrite config with flag values
	confVal := reflect.Indirect(reflect.ValueOf(conf))
	confElem := reflect.ValueOf(conf).Elem()
	for i := 0; i < confVal.NumField(); i++ {
		field := confVal.Type().Field(i)
		fieldType := field.Type
		flagName := field.Tag.Get("flag")

		if fieldType.Kind() == reflect.String {
			flagValue := c.String(flagName)
			if flagValue != "" {
				confElem.Field(i).SetString(flagValue)
			}
		}

		if fieldType.Kind() == reflect.Bool {
			flagValue := c.Bool(flagName)
			if flagValue {
				confElem.Field(i).SetBool(flagValue)
			}
		}
	}

	// Set main token to all servers without token
	for i, server := range conf.Servers {
		if server.Token == "" {
			conf.Servers[i].Token = conf.Token
		}
	}

	// Append plex server from cli flag to list of servers
	plexServer := c.String("plex-server")
	if plexServer != "" && conf.Token != "" {
		conf.Servers = append(conf.Servers, PlexServerConfig{
			BaseURL: plexServer,
			Token:   conf.Token,
		})
	}
}
