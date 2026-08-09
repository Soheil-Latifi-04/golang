package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("loading .env: %v", err)
	}
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable is required", key)
	}
	return value
}

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println("--------------------------------")
	}
}

func main() {
	loadEnv()

	botToken := mustEnv("SLACK_BOT_TOKEN")
	appToken := mustEnv("SLACK_APP_TOKEN")

	bot := slacker.NewClient(botToken, appToken)

	go printCommandEvents(bot.CommandEvents())

	exampleYOB := time.Now().Year() - 30

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "Year of birth calculator",
		Examples:    []string{fmt.Sprintf("my yob is %d", exampleYOB)},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				response.Reply("Invalid year")
				return
			}

			currentYear := time.Now().Year()
			if yob > currentYear || yob < 1900 {
				response.Reply("Please enter a valid year of birth")
				return
			}

			age := currentYear - yob
			response.Reply(fmt.Sprintf("Your age is %d\nYour year of birth is %s", age, year))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
