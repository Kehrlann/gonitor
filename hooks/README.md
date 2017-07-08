# Gonitor / hooks
Here is a collection of useful hooks that you can use with gonitor ; to trigger useful actions on error/recovery
 
## Slack
Use it to send notifications to Slack.

![Slack message](https://raw.githubusercontent.com/Kehrlann/gonitor/master/screenshots/slack.jpg)

- First create a Slack app : https://api.slack.com/slack-apps#creating_apps, with "posting messages" capabilites.
- Then make a copy of slack.sh , say in the same directory as the gonitor binary
- Get the app link for posting, and copy it in your slack.sh at SLACK_URL=...
- Finally, update your gonitor config, usually `gonitor.config.json`, either with a GlobalCommand or a Resource.Command:

```json
    {
      "globalcommand": "/absolute/path/to/your/script/slack.sh",
      "smtp": {...},
      "resources": [...]
    }
```