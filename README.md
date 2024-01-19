# keep-online
[![Super-Linter](https://github.com/rbgayoivoye09/keep-online/actions/workflows/super-linter.yml/badge.svg)](https://github.com/marketplace/actions/super-linter)
自用-周期登入网络认证

## 项目结构
```bash
project_root/
|-- src/                 # 存放源代码
|   |-- main/           # 主要的源代码
|   |-- utils/          # 工具类或辅助功能的源代码
|-- tests/               # 存放测试代码
|-- docs/                # 存放文档
|-- bin/                 # 存放可执行文件
|-- config/              # 存放配置文件
|-- logs/                # 存放日志文件
|-- README.md            # 项目的说明文档
|-- LICENSE              # 项目的许可证
```


## 编译二进制
```bash
 go mod tidy
 go mod vendor
 sh make.sh
```

## 使用

### 修改配置

```bash
cd config 
cp user_demo.yml user.yml # 修改配置
```
### 执行

```bash
$ ./bin/keep-online
Keep online commands

Usage:
  keep-online [command]

Available Commands:
  cfg         Configure keep-online settings
  cmd         Execute a custom command
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  mail        Configure keep-online settings
  ssh         SSH into a remote server

Flags:
  -h, --help   help for keep-online

Use "keep-online [command] --help" for more information about a command.
```
#### 通过邮箱获取VPN认证密码

1. 在配置文件中填写 邮箱名 `name` 和密码 `password` 然后执行命令 `./bin/keep-online mail`


```yml
mail:
  name: "zhangsan@baidu.com" #根据实际情况填写
  password : "zhangsan" # 根据实际情况填写
```

2. 通过命令行参数传入 邮箱名 `name` 和密码 `password`, `./bin/keep-online mail -n zhangsan@baidu.com -p zhangsan`
