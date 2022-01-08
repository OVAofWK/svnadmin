package ini

import (
	"os"
	"path"
	"time"
)

func BackFile(file string) error {
	fileName := path.Base(file)
	currentTime := time.Now().Format("2006-01-02-150405")
	fileName = "backups/" + currentTime + "." + fileName
	info := ReadFile(file)
	f, err := os.Create(fileName)
	Try(err)
	_, err = f.WriteString(info)
	return err
}

func WriteAuthz(file string, conf string, info string) error {
	var authz string

	// 处理需要写的内容
	if conf == "groups" {
		authz = info + "\n" + ReadAuthz(file, "path")
	} else if conf == "path" {
		authz = ReadAuthz(file, "groups") + "\n" + info
	}
	// 写前备份
	err := BackFile(file)
	Try(err)
	f, err := os.Create(file)
	Try(err)
	defer f.Close()
	_, err = f.WriteString(authz)
	return err
}

func WritePasswd(file string, info string) error {
	// 写前备份
	err := BackFile(file)
	Try(err)

	f, err := os.Create(file)
	Try(err)
	defer f.Close()
	_, err = f.WriteString(info)
	return err
}

func WriteConf(CONFIG Config, conf string, info string) error {
	if conf == "groups" || conf == "path" {
		return WriteAuthz(CONFIG.Server.SvnAuthzPath, conf, info)
	} else if conf == "passwd" {
		return WritePasswd(CONFIG.Server.SvnPasswdPath, info)
	}
	return nil
}
