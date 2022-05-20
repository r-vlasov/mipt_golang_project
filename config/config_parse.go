package configuration

import (
    "io/ioutil"
    "path/filepath"
    "gopkg.in/yaml.v2"
    "errors"
)


type DnsTargetConfig struct {
	Name		string	`yaml:"name"`
	Hostname	string	`yaml:"hostname"`
	Repeat		int	`yaml:"repeat"`
	Timeout		int	`yaml:"timeout,omitempty" default: "2"`
	Server		string	`yaml:"server,omitempty"`
}

type PingTargetConfig struct {
	Name		string	`yaml:"name"`	
	Host		string	`yaml:"host"`
	Repeat		int	`yaml:"repeat,omitempty" default: "1"`
	Timeout		int	`yaml:"timeout,omitempty" default: "2"`
}

type HttpTargetConfig struct {
	Name		string	`yaml:"name"`
	Url		string	`yaml:"url"`
	StatusCode	int	`yaml:"code,omitempty" default: "200"`
	Repeat		int	`yaml:"repeat,omitempty" default: "1"`
	Timeout		int	`yaml:"timeout,omitempty" default: "2"`
}

type RipTargetConfig struct {
	Name		string	`yaml:"name"`
	AnnouneAddress	string	`yaml:"announce_address"`
	Timeoit		string	`yaml:"timeout"`
}


type MainConfig struct {
	Version			string 	`yaml:"version"`
	Name			string	`yaml:"name"`
	Description		string	`yaml:"description"`
	NotificationDelay	int  	`yaml:"notify_delay,omitempty" default: 30`
	Targets		struct {
		Dns	[]DnsTargetConfig	`yaml:"dns,omitempty"`
		Ping	[]PingTargetConfig	`yaml:"ping,omitempty"`
		Http	[]HttpTargetConfig	`yaml:"http,omitempty"`
		Rip	[]RipTargetConfig	`yaml:"rip,omitempty"`	// not implemented yet
	} `yaml:"targets"`
}


// config loading function
// 	Arguments: absolute path to config
func (config *MainConfig) Load(absPath string) error {
	// check if path absolute
	if !filepath.IsAbs(absPath) {
		return errors.New("path is not absolute")
	}

	yamlFile, err := ioutil.ReadFile(absPath)

	if err != nil {
		return err
	}

	// deserialization
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return err
	}

	// everything is ok
	return nil
}
