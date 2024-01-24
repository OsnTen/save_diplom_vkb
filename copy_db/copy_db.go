package copy_db

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	pgcommands "github.com/habx/pg-commands"
)

func BackupMySqlDb(host, port, user, password, databaseName string) string {
	var cmd *exec.Cmd

	cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	bytes, err := io.ReadAll(stdout)
	if err != nil {
		log.Println(err)
	}

	now := time.Now().Format("20060102150405")
	var backupPath string

	backupPath = databaseName + "_" + now + ".sql"

	err = os.WriteFile(backupPath, bytes, 0644)

	if err != nil {
		log.Println(err)

	}

	return backupPath
}

func BackupPostgresDb(host, port, user, password, databaseName string) (error, string) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err, ""
	}

	dump, _ := pgcommands.NewDump(&pgcommands.Postgres{
		Host:     host,
		Port:     portInt,
		DB:       databaseName,
		Username: user,
		Password: password,
	})

	dump.Exec(pgcommands.ExecOptions{StreamPrint: false})

	var cmd *exec.Cmd

	cmd = exec.Command("pg_dump", "-h"+host, "-p"+port, "-U"+user, "-W"+password, databaseName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err, ""
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err, ""
	}

	bytes, err := io.ReadAll(stdout)
	if err != nil {
		log.Println(err)
		return err, ""
	}
	// Получить текущую временную метку
	now := time.Now().Format("20060102150405")
	var backupPath string

	backupPath = databaseName + "_" + now + ".sql"

	err = os.WriteFile(backupPath, bytes, 0644)

	if err != nil {
		log.Println(err)
		return err, ""
	}

	return nil, backupPath
}

func BackupReindexer(host, port, databaseName string) string {
	var cmd *exec.Cmd

	fmt.Println("--dsn" + "cproto://" + host + ":" + port + "/" + databaseName)
	now := time.Now().Format("20060102150405")
	var backupPath string

	// Установить имя нашего файла резервного копирования
	backupPath = databaseName + "_" + now + ".rxdump"

	cmd = exec.Command("./reindexer_tool", "-dcproto://"+host+":"+port+"/"+databaseName, "-c\\dump", "-o"+backupPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	bytes, err := io.ReadAll(stdout)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(bytes))

	return backupPath
}

func ReindexerRecover(host, port, databaseName, filename string) {
	var cmd *exec.Cmd

	cmd = exec.Command("reindexer_tool.exe", "--dsn cproto://"+host+":"+port+"/"+databaseName, "--filename"+filename)

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

}

func BackupVM(name string) string {
	var cmd *exec.Cmd

	now := time.Now().Format("20060102150405")
	backupPath := name + "_" + now
	cmds := exec.Command("C:\\Program Files\\Oracle\\VirtualBox\\VBoxManage.exe", "list vms", "--long")
	cmd = exec.Command("C:\\Program Files\\Oracle\\VirtualBox\\VBoxManage.exe", "snapshot"+name, "take"+backupPath)

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
	if err := cmds.Start(); err != nil {
		log.Println(err)
	}

	return backupPath
}
