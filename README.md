[![Go Report Card](https://goreportcard.com/badge/github.com/lowewenzel/go-ack)](https://goreportcard.com/report/github.com/lowewenzel/go-ack)

# Go Acknowledged Bot
_Let your Slack messages be heard and **acknowledged**._

A reacreation of [AcknowledgedBot](https://github.com/lowewenzel/acknowledgedbot), built in [Go](https://github.com/go/go) and for Slack.

### Problem/Proposal

Go-Ack plans to solve the problem of lack of accountability in team chats and groups. Often messages can be missed or read over, but **go-ack** is the Slack tool to minimize that. 

### Vision/Usage

Users can add `@ack` to any group, and when a member wants to create an important message, call ack with the following command.

```
@ack [optional emoji] [required message]
```

The output will be:
```
[...message contents]

----------------
Add a ☑️ emoji to acknowledge!
```

### Implementation
Usage of Open Source Libraries:
- [Slack API by nlopes](https://github.com/nlopes/slack)
- [goslackit by droxey](https://github.com/droxey/goslackit/) 

Features in Development:
- [x] Create Acknowledgement Messages
- [x] Acknowledge messages with emojis
- [ ] Show who has **not** acknowledged yet (in channel)
- [ ] Show `Acknowledge` Button rather than emoji

Future Features:
- [ ] Setup feature per Slack Channel
- [ ] Create messages in App conversation, rather than in group

`License MIT`