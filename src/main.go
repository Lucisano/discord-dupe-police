package main

import (
	"os"

	"github.com/Lucisano/discord-dupe-police/src/bot"
)

func main() {
	bot.BotToken = os.Getenv("DUPE_POLICE_DISCORD_BOT_TOKEN")
	bot.Run() // call the run function of bot/bot.go
}
