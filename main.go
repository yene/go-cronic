package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/naoina/toml"
	mail "github.com/xhit/go-simple-mail"
)

var defaultConfigPath = "~/.config/cronic/cronic.conf"

func main() {
	configPathPtr := flag.String("c", defaultConfigPath, "Path to config")
	flag.Parse()
	configPath := *configPathPtr

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("Require a command")
	}

	_ = godotenv.Load()
	config := TomlConfig{}
	config.Mail.Template = defaultMailTemplate()

	if fileExists(expandTilde(configPath)) {
		dat, err := ioutil.ReadFile(expandTilde(configPath))
		if err != nil {
			log.Fatalln(err)
		}
		toml.Unmarshal(dat, &config)
		// log.Println("loading config file", expandTilde(configPath))
	}

	config = LoadFromENV(config)

	var outbuf, errbuf bytes.Buffer
	var exitCode int = 0
	cmd := exec.Command(args[0], args[1:]...)
	commandString := strings.Join(args, " ")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	// https://stackoverflow.com/a/55055100/279890
	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode = exitError.ExitCode()
	} else {
		log.Println("no error")
		if config.Mail.Sendstdout == false {
			os.Exit(0)
		}
	}

	data := TemplateVariables{
		Command:        commandString,
		ResultCode:     exitCode,
		ErrorOutput:    errbuf.String(),
		StandardOutput: outbuf.String(),
		TraceError:     "",
	}

	var tpl bytes.Buffer
	t := template.Must(template.New("html-tmpl").Parse(config.Mail.Template))
	err = t.Execute(&tpl, data)
	if err != nil {
		log.Fatalln(err)
	}
	htmlBody := tpl.String()

	// https://github.com/xhit/go-simple-mail
	server := mail.NewSMTPClient()
	server.Host = config.Smtp.Host
	server.Port = config.Smtp.Port
	server.Username = config.Smtp.Username
	server.Password = config.Smtp.Password
	server.Encryption = convertEncryption(config.Smtp.Encryption)
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	sus := commandString
	if len(sus) > 25 {
		sus = sus[0:22] + "..."
	}
	subject := "Cronic error for: "
	if exitCode == 0 {
		subject = "Cronic success for: "
	}
	subject = subject + sus

	if config.Mail.Subject != "" {
		subject = config.Mail.Subject
	}

	email := mail.NewMSG()
	email.SetFrom("Cronic <" + config.Mail.Sender + ">").
		AddTo(config.Mail.Receiver).
		SetSubject(subject)

	email.SetBody(mail.TextPlain, htmlBody)
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}
	os.Exit(exitCode)
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err == nil {
			return filepath.Join(usr.HomeDir, path[2:])
		}
	}
	return path
}

func convertEncryption(enc string) mail.Encryption {
	if enc == "SSL" {
		return mail.EncryptionSSL
	} else if enc == "TLS" {
		return mail.EncryptionTLS
	} else {
		return mail.EncryptionNone
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
