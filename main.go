package main

import (
	"fmt"
	"os"

	"github.com/lowewenzel/goack/goacklib"

	"github.com/go-redis/redis"
	_ "github.com/joho/godotenv/autoload"
)

// main is our entrypoint, where the application initializes the Slackbot.
// DO NOT EDIT THIS FUNCTION. This is a fully complete implementation.
func main() {
	slackIt()
}

// RedisNewClient Start the Redis server
func RedisNewClient() *redis.Client {
	fmt.Println("Connecting to " + os.Getenv("REDIS_URL"))
	Client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PW"), // no password set
		DB:       0,                     // use default DB
	})

	fmt.Println("Connected to Client!")
	return Client
}

func slackIt() {
	redisClient := RedisNewClient()
	botToken := os.Getenv("BOT_OAUTH_ACCESS_TOKEN")
	slackClient := goacklib.SlackLib{}
	slackClient.CreateSlackClient(botToken)
	go goacklib.RunServer(slackClient.RTM, redisClient)
	goacklib.RespondToEvents(slackClient.RTM, redisClient)
}
