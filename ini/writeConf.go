package ini

import (
	"os"
)

func WriteAuthz(file string, conf string, info string) error {
	var authz string
	// 处理需要写的内容
	if conf == "groups" {
		authz = info + "\n" + ReadAuthz(file, "path")
	} else if conf == "path" {
		authz = ReadAuthz(file, "groups") + "\n" + info
	}
	f, err := os.Create(file)
	Try(err)
	defer f.Close()
	_, err = f.WriteString(authz)
	return err
}

func WritePasswd(file string, info string) error {
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
