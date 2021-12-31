package ini

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Listen        string `yaml:"listen"`
		SvnAuthzPath  string `yaml:"svnAuthzPath"`
		SvnPasswdPath string `yaml:"svnPasswdPath"`
	} `yaml:"server"`

	Admin struct {
		User   string `yaml:"user"`
		Passwd string `yaml:"passwd"`
	} `yaml:"admin"`
	Web struct {
		Title string `yaml:"title"`
	}
}

func ReadConfYaml(confPath string) Config {
	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("read %s failed", confPath)
		panic(err)
	}
	return config
}
