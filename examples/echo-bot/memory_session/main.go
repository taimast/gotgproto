package main

import (
	"fmt"
	"log"

	"github.com/taimast/gotgproto"
	"github.com/taimast/gotgproto/dispatcher"
	"github.com/taimast/gotgproto/dispatcher/handlers"
	"github.com/taimast/gotgproto/dispatcher/handlers/filters"
	"github.com/taimast/gotgproto/ext"
	"github.com/taimast/gotgproto/sessionMaker"
	"github.com/gotd/td/tg"
)

func main() {
	client, err := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		123456,
		// Get ApiHash from https://my.telegram.org/apps
		"API_HASH_HERE",
		// ClientType, as we defined above
		gotgproto.ClientTypeBot("BOT_TOKEN_HERE"),
		// Optional parameters of client
		&gotgproto.ClientOpts{
			InMemory: true,
			Session:  sessionMaker.SimpleSession(),
		},
	)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}

	clientDispatcher := client.Dispatcher

	// Command Handler for /start
	clientDispatcher.AddHandler(handlers.NewCommand("start", start))
	// Callback Query Handler with prefix filter for recieving specific query
	clientDispatcher.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("cb_"), buttonCallback))
	// This Message Handler will call our echo function on new messages
	clientDispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, echo), 1)

	fmt.Printf("client (@%s) has been started...\n", client.Self.Username)

	err = client.Idle()
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}
}

// callback function for /start command
func start(ctx *ext.Context, update *ext.Update) error {
	user := update.EffectiveUser()
	_, _ = ctx.Reply(update, ext.ReplyTextString(fmt.Sprintf("Hello %s, I am @%s and will repeat all your messages.\nI was made using gotd and gotgproto.", user.FirstName, ctx.Self.Username)), &ext.ReplyOpts{
		Markup: &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonURL{
							Text: "gotd/td",
							URL:  "https://github.com/gotd/td",
						},
						&tg.KeyboardButtonURL{
							Text: "gotgproto",
							URL:  "https://github.com/taimast/gotgproto",
						},
					},
				},
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "Click Here",
							Data: []byte("cb_pressed"),
						},
					},
				},
			},
		},
	})
	// End dispatcher groups so that bot doesn't echo /start command usage
	return dispatcher.EndGroups
}

func buttonCallback(ctx *ext.Context, update *ext.Update) error {
	query := update.CallbackQuery
	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		Alert:   true,
		QueryID: query.QueryID,
		Message: "This is an example bot!",
	})
	return nil
}

func echo(ctx *ext.Context, update *ext.Update) error {
	msg := update.EffectiveMessage
	_, err := ctx.Reply(update, ext.ReplyTextString(msg.Text), nil)
	return err
}
