package ini

import (
	"bufio"
	"os"
	"strings"
)

func Try(err error) {
	if err != nil {
		panic(err)
	}
}

//读authz配置，返回一个 配置项:配置内容 的map
func ReadAuthz(file string, conf string) string {
	f, err := os.Open(file)
	Try(err)
	defer f.Close()
	bfread := bufio.NewReader(f)

	list := make(map[string]string)
	option := ""
	//开始读取文件
	for {
		line, _, err := bfread.ReadLine()
		if err != nil {
			break
		}
		if len(line) > 1 {
			sline := string(line)
			if strings.Contains(sline, "[groups]") {
				option = "groups"
				list[option] = string(sline)
			} else if strings.Contains(sline, "[") && strings.Contains(sline, "]") {
				option = "path"
				list[option] = list[option] + "\n" + sline
			} else {
				list[option] = list[option] + "\n" + sline
			}
		}
	}
	Try(err)
	return list[conf]

}

//读取passwd需求配置
func ReadPasswd(file string) string {
	f, err := os.Open(file)
	Try(err)
	defer f.Close()
	bfread := bufio.NewReader(f)
	passwd := ""
	for {
		line, _, err := bfread.ReadLine()
		sline := string(line)
		if err != nil {
			break
		}
		if len(line) > 1 {
			passwd = passwd + "\n" + sline
		}
	}
	return passwd
}

func ReadConf(CONFIG Config, conf string) string {
	var info string
	if conf == "groups" || conf == "path" {
		info = ReadAuthz(CONFIG.Server.SvnAuthzPath, conf)
	} else if conf == "passwd" {
		info = ReadPasswd(CONFIG.Server.SvnPasswdPath)
	}

	return info
}
