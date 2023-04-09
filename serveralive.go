package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

// Bot parameters
var (
	GuildID        = ""
	ChannelID      = ""
	BotToken       = ""
	RemoveCommands = true
)

var s *discordgo.Session

var (
	dmPermission = true

	commands = []*discordgo.ApplicationCommand{
		{
			Name:         "ip",
			Description:  "Get server public IP",
			DMPermission: &dmPermission,
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ip": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: cuteIpResponse(),
				},
			})
		},
	}
)

func getIPAddress() (string, error) {
	resp, err := http.Get("http://ifconfig.me/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert the response body to a string and return it
	return string(body), nil
}

func cuteIpResponse() string {
	ip, err := getIPAddress()
	if err != nil {
		return fmt.Sprintf("Cannot determine server IP:\n```\n%s\n```", err)
	} else {
		return fmt.Sprintf("My IP is `%s`", ip)
	}
}

func init() {
	// Init env
	GuildID = os.Getenv("DC_GUILD_ID")
	ChannelID = os.Getenv("DC_CHANNEL_ID")
	BotToken = os.Getenv("DC_TOKEN")
	if BotToken == "" {
		log.Fatalln("Bot token cannot be empty")
	}
}

func main() {
	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		if ChannelID != "" {
			_, err := s.ChannelMessageSend(ChannelID, "Server is up\n"+cuteIpResponse())
			if err != nil {
				log.Println("Cannot sent start up message: ", err)
			}
		}
	})
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
