package goack

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/nlopes/slack"
)

/*
   TODO: Change @BOT_NAME to the same thing you entered when creating your Slack application.
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = "Type in `@AcknowledgedBot <emoji> <message>` to make a message!"

/*CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
  initiating the socket connection and returning the client.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*RespondToEvents waits for messages on the Slack client's incomingEvents channel,
  and sends a response when it detects the bot has been tagged in a message with @<botTag>.

  EDIT THIS FUNCTION IN THE SPACE INDICATED ONLY!
*/
func RespondToEvents(slackClient *slack.RTM, redisClient *redis.Client) {
	for msg := range slackClient.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)
			// TODO: Make your bot do more than respond to a help command. See notes below.
			// Make changes below this line and add additional funcs to support your bot's functionality.
			// sendHelp is provided as a simple example. Your team may want to call a free external API
			// in a function called sendResponse that you'd create below the definition of sendHelp,
			// and call in this context to ensure execution when the bot receives an event.

			// START SLACKBOT CUSTOM CODE
			// ===============================================================

			switch command := strings.Fields(message)[0]; command {
			case "help":
				sendHelp(slackClient, message, ev.Channel)
			// case "users":
			// sendUsers(slackClient, message, ev.Channel)
			default:
				sendResponse(slackClient, message, ev.Channel, redisClient)
			}

			// ===============================================================
			// END SLACKBOT CUSTOM CODE
		default:
		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendHelp is a working help message, for reference.
func getUsers(slackClient *slack.RTM, message, slackChannel string, redisClient *redis.Client) (res string) {
	params := slack.GetUsersInConversationParameters{ChannelID: slackChannel, Cursor: "", Limit: 0}
	users, _, err := slackClient.GetUsersInConversation(&params)
	if err != nil {
		fmt.Println(err)
	}
	for _, user := range users {
		redisClient.HSet(slackChannel, user, "YEET")
		res += ("<@" + string(user) + "> ")
	}
	return res
}

func sendResponse(slackClient *slack.RTM, message, slackChannel string, redisClient *redis.Client) {
	// command := strings.ToLower(message)
	emoji := strings.Fields(message)[0]
	if !strings.Contains(emoji, ":") {
		emoji = ":heavy_check_mark:"
	}

	footerText := "Not Acknowledged:\n" + getUsers(slackClient, message, slackChannel, redisClient)

	attachmentAction := slack.AttachmentAction{Name: "Ack", Text: "Acknowledge", Type: "button"}
	attachment := slack.Attachment{
		Pretext:    message,
		Text:       footerText,
		Fallback:   "Update Slack to view this message!",
		Actions:    []slack.AttachmentAction{attachmentAction},
		CallbackID: "ack_msg",
		Color:      "#479ACC",
	}

	// newMessage := message + footerText
	slackClient.PostMessage(slackChannel, slack.MsgOptionAttachments(attachment))
	// slackClient.SendMessage(slackClient.NewOutgoingMessage(newMessage, slackChannel))
}

func acknowledgeCallback(c echo.Context, slackClient *slack.RTM, redisClient *redis.Client) {
	jsonString := c.FormValue("payload")
	var Callback slack.InteractionCallback
	json.Unmarshal([]byte(jsonString), &Callback)
	// slackClient.PostMessage(Callback.Channel.ID, slack.MsgOptionText("Acknowledged!", false))

	channelUsers := redisClient.HGetAll(Callback.Channel.ID)
	// fmt.Println(channelUsers.Val())
	finalMessage := "Not Acknowledged:\n"

	for i := range channelUsers.Val() {
		if i != Callback.User.ID {
			finalMessage += ("<@" + string(i) + "> ")
		}
	}

	newAttachment := slack.Attachment{
		Pretext:    Callback.OriginalMessage.Attachments[0].Pretext,
		Text:       finalMessage,
		Fallback:   "Update Slack to view this message!",
		Actions:    Callback.OriginalMessage.Attachments[0].Actions,
		CallbackID: "ack_msg",
		Color:      "#479ACC",
	}

	newMessage := slack.MsgOptionAttachments(newAttachment)

	// slackClient.PostMessage(Callback.Channel.ID, newMessage)
	slackClient.UpdateMessage(Callback.Channel.ID, Callback.OriginalMessage.Timestamp, newMessage)
	// slackClient.DeleteMessage(Callback.Channel.ID, Callback.OriginalMessage.Timestamp)

	// return newMessage
}

// func deleteAllMessages(c echo.Context, slackClient *slack.RTM) {
// 	slackClient.GetChannelHistory()
// }

// RunServer starts the Slack Bot server for attachment callbacks
func RunServer(slackClient *slack.RTM, redisClient *redis.Client) {
	port := ":" + os.Getenv("PORT")
	e := echo.New()

	// Callback for Acknowledge Button
	e.POST("/api", func(c echo.Context) error {
		acknowledgeCallback(c, slackClient, redisClient)
		return nil
	})
	e.Logger.Fatal(e.Start(port))
}
