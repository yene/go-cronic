package main

import (
	"bytes"
	"crypto/tls"
	_ "embed"
	"errors"
	"flag"
	"html/template"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/naoina/toml"
	mail "github.com/xhit/go-simple-mail/v2"
)

//go:embed template.txt
var defaultMailTemplate string

var defaultConfigPath = "~/.config/cronic/cronic.conf"
var validate = false

func main() {
	configPathPtr := flag.String("c", defaultConfigPath, "Path to config")
	flag.Parse()
	configPath := *configPathPtr

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("Require a command")
	}

	if args[0] == "validate" {
		validate = true
	}

	_ = godotenv.Load()
	config := TomlConfig{}

	config.Mail.Template = defaultMailTemplate

	if fileExists(expandTilde(configPath)) {
		dat, err := os.ReadFile(expandTilde(configPath))
		if err != nil {
			log.Fatalln(err)
		}
		toml.Unmarshal(dat, &config)
	}

	config = LoadFromENV(config)

	var outbuf, errbuf bytes.Buffer
	var exitCode int = 0
	cmd := exec.Command(args[0], args[1:]...)
	commandString := strings.Join(args, " ")
	if validate {
		commandString = "echo validate"
	}

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	// https://stackoverflow.com/a/62647366/279890

	var ee *exec.ExitError
	var pe *os.PathError
	if errors.As(err, &ee) {
		exitCode = ee.ExitCode()
	} else if errors.As(err, &pe) {
		// "127 no such file ...", "126 permission denied" etc.
		exitCode = 1
	} else if err != nil {
		// "127 executable file not found in $PATH"
		exitCode = 1
	}

	if exitCode == 0 {
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
	if validate {
		log.Println("Succesfully connected to", config.Smtp.Host, config.Smtp.Port)
		os.Exit(0)
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
