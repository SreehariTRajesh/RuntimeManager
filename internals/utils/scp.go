package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func TransferFilesSFTP(container_id string, dest_ip string, dest_user string, pwd string) error {
	config := &ssh.ClientConfig{
		User: dest_user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:22", dest_ip)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("error while connecting to remote host: %w", err)
	}
	defer conn.Close()
	source_file_path := fmt.Sprintf("/tmp/%s.tar.gz", container_id)
	destination_file_path := fmt.Sprintf("/tmp/%s.tar.gz", container_id)
	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("error while creating sftp client: %w", err)
	}
	defer client.Close()

	source_file, err := os.Open(source_file_path)
	if err != nil {
		return fmt.Errorf("error while opening source file: %w", err)
	}
	defer source_file.Close()

	destination_file, err := client.Create(destination_file_path)
	if err != nil {
		return fmt.Errorf("error while opening destination file: %w", err)
	}
	defer destination_file.Close()

	_, err = io.Copy(destination_file, source_file)

	if err != nil {
		return fmt.Errorf("error while copying file from source to destination; %w", err)
	}
	return nil
}
