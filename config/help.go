package config

import "fmt"

func PrintUsage() {
    fmt.Printf(`Gonitor gets its configuration from a JSON file. You can get an example template by running.

	gonitor --example

CONFIG DETAILS :
----------------
- Resources
    Is an array describing which websites to monitor, and the details of how to
    monitor them. Each resource looks like this :

    - url                   the url of the website to monitor, including
                            scheme, e.g. "http://www.example.com"

    - intervalInSeconds     the interval at which this url is polled, e.g. 30
                            for "once every 30 seconds"

    - timeoutInSeconds      the acceptable timeout when polling, e.g. 1. If an
                            HTTP GET times out, it is considered a failure,
                            with a status code "0"

    - numberOfTries         the number of successive calls to consider to know
                            whether a service is up or not

    - failureThreshold      the number of failures to trigger a failure

    - command [optional]    a command to execute any arbitrary process when a
                            failure or recovery is detected.
                            If "GlobalCommand" (see below) is defined, this
                            value takes precedence over GlobalCommand: only
                            this command will be executed.
                            Examples can be found at:
                            https://github.com/Kehrlann/gonitor/tree/master/hooks


- GlobalCommand [optional]
    It is a convenient way to execute arbitrary commands when a failure or
    recovery is detected.

    The command gets four arguments passed in :
        - The message type as string, either "RECOVERY" or "FAILURE"
        - The datetime of the message
        - The last HTTP codes
        - The resource URL

    So putting "echo" as a command would print :
        'RECOVERY 2017-07-14T13:45:03 [200 0 0 0 0] http://www.example.com/'

    You can find useful scripts at : 
    https://github.com/Kehrlann/gonitor/tree/master/hooks


- SMTP [optional]
    Is the SMTP config for sending alert e-mails. If this is not setup in the 
    config, no mail will be sent when a failure / recovery occurs.

    - host              the hostname, e.g. "smtp.example.com"

    - port              the SMTP port to connect to, as an int, e.g. 25

    - username          the username for connecting to the server

    - password          the password for connecting to the server

    - fromAddress       the sender's address used in the "From:" header, e.g.
                        "gonitor@example.com"

    - fromName          the sender's name to be used in the "From:" header, 
                        e.g. "Gonitor". If fromAddress is "gonitor@example.com"
                        the "From:" will be "Gonitor <gonitor@example.com>"

    - to                an array containing the list of recipients, e.g. 
                        ["recipient@example.com", "admin@example.com"]


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

`)
}

func PrintExampleConfig() {
    fmt.Println(
        `{
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
}`)
}