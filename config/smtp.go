package config

import "fmt"

// Smtp is the config par that describes how to talk to an SMTP server and send e-mails
type Smtp struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	To          []string
}

// IsValid tells you whether you can trust this Smtp config to send an e-mail
func (smtp *Smtp) IsValid() bool {
	return smtp.Host != "" && smtp.Port != 0 && len(smtp.To) > 0 && smtp.FromAddress != ""
}

// FormatFromHeader creates the From header used in SMTP messages
func (smtp *Smtp) FormatFromHeader() string {
	return fmt.Sprintf("%v <%v>", smtp.FromName, smtp.FromAddress)
}
