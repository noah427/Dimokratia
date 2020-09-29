package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var client *discordgo.Session

func main() {
	godotenv.Load()
	initDiscord()
}

func initDiscord() {
	client, _ = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	err := client.Open()
	if err != nil {
		log.Fatalln("Invalid Token")
	}

	log.Printf("Now running | Logged in as %s\n", client.State.User.Username)


	//exit
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Now exiting")

	client.Close()
}
