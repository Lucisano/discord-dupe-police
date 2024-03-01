package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)
	//discord.Identify.Intents = discordgo.IntentsAll
	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "http"):
		returnChan := make(chan *discordgo.Message)
		go searchForLink(discord, message, returnChan)
		historicalMessage := <-returnChan
		if historicalMessage != nil {
			discord.ChannelMessageSendReply(historicalMessage.ChannelID, "It looks like you're a wally, here's the same link you've just posted, but earlier!", &discordgo.MessageReference{
				MessageID: historicalMessage.ID,
			})
		}
	}

}

func searchForLink(discord *discordgo.Session, message *discordgo.MessageCreate, returnChan chan *discordgo.Message) {
	messages, err := discord.ChannelMessages(message.ChannelID, 100, "", "", "")
	if err != nil {
		fmt.Println("Error retrieving message history:", err)
		return
	}

	checkUrl := stripUrl(message.Content)

	for _, msg := range messages {
		if msg.ID == message.ID {
			continue
		}
		historicalMessageUrl := stripUrl(msg.Content)
		if historicalMessageUrl == "" {
			continue
		}
		if strings.Contains(historicalMessageUrl, checkUrl) {
			returnChan <- msg
			return
		}
	}
	returnChan <- nil
}

func stripUrl(inputText string) string {
	urlRegex := regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]*[-A-Za-z0-9+&@#/%=~_|]`)

	// Find the first match of the URL pattern in the text
	url := urlRegex.FindString(inputText)

	// Print the extracted URL
	return url
}
