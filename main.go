package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {

	botToken := flag.String("token", "", "Bot token")
	flag.Parse()

	if *botToken == "" {
		fmt.Println("BOT_TOKEN is not set")
		os.Exit(1)
	}
	fmt.Println("Bot token:'", *botToken, "'")
	// create new bot

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(*botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	// Initialize done chan
	done := make(chan struct{}, 1)

	// Get bot user
	botUser, err := bot.GetMe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Bot user: %+v\n", botUser)

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	// Create bot handler
	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()

	// Handle any command
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)
		fmt.Println("Got update: ", update)
		//time.Sleep(time.Second * 15) // Simulate long process time
		// Send a message with provided message entities
		_, _ = bot.SendMessage(tu.MessageWithEntities(chatID,
			//tu.Entity("Hi").Bold(), tu.Entity(" "), tu.Entity("There").Italic().Spoiler(), tu.Entity("\n"),
			tu.Entity("Hi "), tu.Entity(update.Message.From.Username), tu.Entity("\n"),
			tu.Entity("JIRA Tasks: "), tu.Entity("\n"),
			tu.Entity("TASK-5436").TextLink("https://example.com").Italic(), tu.Entity("\n"),
			tu.Entity("TASK-5436").TextLink("https://example.com").Italic(), tu.Entity("\n"),
			tu.Entity("TASK-5436").TextLink("https://example.com").Italic(), tu.Entity("\n"),
			tu.Entity("TASK-5436").TextLink("https://example.com").Italic(), tu.Entity("\n"),
		))

	}, th.CommandEqual("task"))

	// // Handle any message
	// bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
	// 	// Get chat ID from the message
	// 	chatID := tu.ID(message.Chat.ID)

	// 	// Copy sent messages back to the user
	// 	_, _ = bot.CopyMessage(
	// 		tu.CopyMessage(chatID, chatID, message.MessageID),
	// 	)
	// })

	// Handle stop signal (Ctrl+C)
	go func() {
		// Wait for stop signal
		<-sigs

		fmt.Println("Stopping...")

		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		//defer cancel()

		bot.StopLongPolling()
		fmt.Println("Long polling done")

		//bh.StopWithContext(ctx)
		//fmt.Println("Bot handler done")

		// Notify that stop is done
		done <- struct{}{}
	}()

	// Start handling updates
	bh.Start()

	fmt.Println("Handling updates...")

	// Wait for the stop process to be completed
	<-done
	fmt.Println("Done")
}
