package ini

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Listen        string `yaml:"listen"`
		SvnAuthzPath  string `yaml:"svnAuthzPath"`
		SvnPasswdPath string `yaml:"svnPasswdPath"`
	} `yaml:"server"`

	Admin struct {
		User    string `yaml:"user"`
		Passwd  string `yaml:"passwd"`
		LogPath string `yaml:"logPath"`
		UseLog  bool   `yaml:"useLog"`
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
func isAxist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
