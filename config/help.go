package config

import "fmt"

var EXAMPLE_CONFIG = `{
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
}`

func PrintExampleConfig() {
	fmt.Printf(`Gonitor gets its configuration from a JSON file. Below is a starter template. The differe
You can either replace the SMTP config with your own, or remove it entirely if you don't want e-mail notifications.

CONFIG DETAILS :
----------------
- GlobalCommand
	Is a convenient way to execute arbitrary commands when a message (Failure/Recovery) is emitted by Gonitor.
	You get four arguments :
		- The message type as string, either "RECOVERY" or "FAILURE"
		- The datetime of the message
		- The last HTTP codes
	 	- The resource URL

	So putting "echo" as a command would print :
		'RECOVERY ....'
		TODO
		TODO
		TODO
		TODO
		TODO
		TODO
		TODO
		TODO
		TODO

TEMPLATE :
----------
%v`, EXAMPLE_CONFIG)
}