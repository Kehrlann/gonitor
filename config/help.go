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
      "command"           : "/home/user/scripts/database.sh"
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
        'RECOVERY 2017-07-14T13:45:03 [200 0 0 0 0] http://www.example.com/'

- SMTP
    Is the SMTP config for sending alert e-mails. If this is not setup in the json, the config is still valid, but no
    mail will be sent when a failure / recovery occurs.
    - host:            the hostname, e.g. "smtp.example.com"
    - port:            the SMTP port to connect to, as an int, e.g. 25
    - username:        the username for connecting to the server
    - password:        the password for connecting to the server
    - fromAddress:     the sender's address used in the "From:" header, e.g. "gonitor@example.com"
    - fromName:        the sender's name to be used in the "From:" header, e.g. "Gonitor"
                       using both examples above, the "From:" would be "Gonitor <gonitor@example.com>"
    - to:              an array containing the list of recipients, e.g. ["recipient@example.com", "admin@example.com"]

- Resources
    Is an array describing which websites to monitor, and the details of how to monitor them.
    - url                    the url of the website to monitor, including scheme, e.g. "http://www.example.com"
    - intervalInSeconds      the interval at which this url is polled, e.g. 30 for "once every 30 seconds"
    - timeoutInSeconds       the acceptable timeout when polling, e.g. 1.
                             if an HTTP GET times out, it is considered a failure with a 0 response code
    - numberOfTries          the number of successive calls to consider to know whether a service is up or not
    - failureThreshold       the number of failures to trigger a failure
    - command                a specific command to execute, that overrides the GlobalCommand
                             for additional details see GlobalCommand

USE CASE :
----------
- number of tries : 5
- failure threshold : 3

    Last 5 http codes :     | Result :
    ------------------------+-----------
    [200 200 200 200 404]   | nothing
    ------------------------+-----------
    [200 404 404 200 404]   | failure
    ------------------------+-----------
    [200 200 404 404 404]   | failure
    ------------------------+-----------
    [200 200 200 200 200]   | recovery


TEMPLATE :
----------
%v`, EXAMPLE_CONFIG)
}