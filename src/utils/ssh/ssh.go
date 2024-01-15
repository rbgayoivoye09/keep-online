package ssh

import (
	"fmt"
	"time"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"

	sys_ssh "golang.org/x/crypto/ssh"
)

// RemoteFileContent 获取远程文件的内容(VPN密码)
func RemoteFileContent(username, password, host, filePath string, port int) (string, error) {
	// Logger.Sugar().Infof("SSH into the remote server: %s@%s:%d %s", username, host, port, filePath)
	// SSH配置
	config := &sys_ssh.ClientConfig{
		User: username,
		Auth: []sys_ssh.AuthMethod{
			sys_ssh.Password(password),
		},
		HostKeyCallback: sys_ssh.InsecureIgnoreHostKey(),
	}

	t := time.Now()
	// 连接SSH服务器
	client, err := sys_ssh.Dial("tcp", host+":"+fmt.Sprint(port), config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	Logger.Sugar().Info("连接SSH服务器", time.Since(t))

	t2 := time.Now()
	// 打开一个新的会话
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	Logger.Sugar().Info("打开一个新的会话", time.Since(t2))

	// 执行远程命令
	t3 := time.Now()
	cmd := "cat " + filePath
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}
	Logger.Sugar().Info("执行远程命令", time.Since(t3))

	return string(output), nil
}
