# svnadmin


# 前言
这是最基础的版本，只实现了修改SVN配置的最基本功能，没有其他功能

# 使用方式
svnadmin使用很简单，只有一个二进制文件和一个yaml配置文件


## 程序内容
配置文件默认存放在svnadmin二进制文件目录的conf下。现在的目录结构是这样的:
```
├── conf
│   └── svnconf.yaml
└── svnadmin
```
## 启动
```
./svnadmin
```

# 配置文件解释
```yaml
server:
  listen: 0.0.0.0:88         #监听端口 也可以使用 :88的方式
  svnAuthzPath: "conf/authz" #SVN的 authz配置文件的路径
  svnPasswdPath: conf/passwd #SVN的 passwd配置文件路径
web:
  title: "svn网页版配置"      #web页面的标题
admin:
  user: svnadmin             #登录账号
  passwd: 123456             #密码
```

