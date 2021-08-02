package main

import (
	"log"
	"os"
	"strconv"
)

func LoadFromENV(c TomlConfig) TomlConfig {
	prefix := "CRONIC_MAIL_"
	c.Mail.Sender = lookupEnv(prefix+"SENDER", c.Mail.Sender)
	c.Mail.Receiver = lookupEnv(prefix+"RECEIVER", c.Mail.Receiver)
	c.Mail.Sendstdout = lookupBool(prefix+"SENDSTDOUT", c.Mail.Sendstdout)
	c.Mail.Subject = lookupEnv(prefix+"SUBJECT", c.Mail.Subject)
	c.Mail.Template = lookupEnv(prefix+"TEMPLATE", c.Mail.Template)

	prefix = "CRONIC_SMTP_"
	c.Smtp.Host = lookupEnv(prefix+"HOST", c.Smtp.Host)
	c.Smtp.Port = lookupInt(prefix+"PORT", c.Smtp.Port)
	c.Smtp.Encryption = lookupEnv(prefix+"ENCRYPTION", c.Smtp.Encryption)
	c.Smtp.Authentication = lookupEnv(prefix+"AUTHENTICATION", c.Smtp.Authentication)
	c.Smtp.Username = lookupEnv(prefix+"USERNAME", c.Smtp.Username)
	c.Smtp.Password = lookupEnv(prefix+"PASSWORD", c.Smtp.Password)
	return c
}

type TemplateVariables struct {
	Command        string
	ResultCode     int
	ErrorOutput    string
	StandardOutput string
	TraceError     string
}

type TomlConfig struct {
	Smtp ConfigSMTP `toml:"smtp"`
	Mail ConfigMail `toml:"mail"`
}

type ConfigSMTP struct {
	Host           string `toml:"host"`
	Port           int    `toml:"port"`
	Encryption     string `toml:"encryption"`
	Authentication string `toml:"authentication"`
	Username       string `toml:"username"`
	Password       string `toml:"password"`
}

type ConfigMail struct {
	Sender     string `toml:"sender"`
	Receiver   string `toml:"receiver"`
	Sendstdout bool   `toml:"sendstdout"`
	Subject    string `toml:"subject"`
	Template   string `toml:"template"`
}

func lookupEnv(key string, current string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return current
	} else {
		return val
	}
}

func lookupBool(key string, current bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return current
	} else {
		if val == "true" || val == "True" || val == "1" {
			return true
		} else if val == "false" || val == "False" || val == "0" {
			return false
		} else {
			log.Printf("%q is not a bool.\n", key)
			return current
		}
	}
}

func lookupInt(key string, current int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return current
	} else {
		if num, err := strconv.Atoi(val); err == nil {
			return num
		} else {
			log.Printf("%q is not a number.\n", key)
			return current
		}
	}
}
