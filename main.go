package main

import (
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
	Client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return Client
}

func slackIt() {
	redisClient := RedisNewClient()
	botToken := os.Getenv("BOT_OAUTH_ACCESS_TOKEN")
	slackClient := goacklib.CreateSlackClient(botToken)
	go goacklib.RunServer(slackClient, redisClient)
	goacklib.RespondToEvents(slackClient, redisClient)
}
