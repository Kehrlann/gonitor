package config

import "fmt"

// NoDefaultConfigError represents an error thrown when the user hasn't specified a config, and that config wasn't found
type NoDefaultConfigError struct {
	HelpMessage string
}

// Error is the error message for NoDefautConfigError
func (n NoDefaultConfigError) Error() string {
	return fmt.Sprintf("No default config found at %v", DEFAULT_CONFIG_PATH)
}

func NewDefaultConfigError() *NoDefaultConfigError {
	return &NoDefaultConfigError{HelpMessage: fmt.Sprintf(`
It seems you didn't specify a config file. Gonitor attempted to load a default config file, located at  :
	%v

No such file was found. Please create one, or specify an existing config with the -config flag. If you wish to create
one, here is a starter template. You can either replace the SMTP config with your own, or remove it entirely if you
don't want e-mail notifications.

{
  "globalcommand" : "/home/user/scripts/slack.sh",
  "smtp"    :
  {
    "host"        : "smtp.example.com",
    "port"        : 25,
    "username"    : "user",
    "password"    : "password123",
    "fromaddress" : "address@example.com",
    "fromname"    : "Mr Example",
    "to"          : ["recipient@example.com", "admin@example.com"]
  },
  "resources" :
  [
	{
	  "url"               : "http://www.example.com",
	  "intervalInSeconds" : 60,
	  "timeoutInSeconds"  : 2,
	  "numberOfTries"     : 10,
	  "failureThreshold"  : 3
	},
	{
	  "url"               : "http://www.example.test",
	  "intervalInSeconds" : 120,
	  "timeoutInSeconds"  : 10,
	  "numberOfTries"     : 10,
	  "failureThreshold"  : 10,
	  "command"			  : "/home/user/scripts/database.sh"
	}
  ]
}`, DEFAULT_CONFIG_PATH)}
}
