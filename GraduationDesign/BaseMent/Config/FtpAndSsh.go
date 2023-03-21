package Config

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"time"
)

func InitFtp() {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(Conf.SFtp.Password))
	clientConfig = &ssh.ClientConfig{
		User:            Conf.SFtp.User,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}
	// connet to ssh
	port, _ := strconv.Atoi(Conf.SFtp.Port)
	addr = fmt.Sprintf("%s:%d", Conf.SFtp.Host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		log.Println("ssh连接失败", err)
	}
	// create sftp client
	if SftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Println("sftp连接失败", err)
	}
	log.Println("sftpInfo:", SftpClient)
}
