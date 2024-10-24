package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func main() {

	var (
		Token   = "Bot " + os.Getenv("DISCORD_TOKEN")
		BotName = "<@" + os.Getenv("DISCORD_CLIENT_ID") + ">"
	)

	go listenPort()
	
	fmt.Println(Token)
	fmt.Println(BotName)

	discord, err := discordgo.New(Token)
	
	if err != nil {
		fmt.Println("ログインに失敗しました")
		fmt.Println(err)
	}

	//イベントハンドラを追加
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})
	
	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}
	// 直近の関数（main）の最後に実行される
	defer discord.Close()
	
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
	
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	clientId := os.Getenv("DISCORD_CLIENT_ID")

	u := m.Author
	
	if u.ID == clientId {
		return
	}

	channel, _ := 	s.State.Channel(m.ChannelID)
	if channel.Name != "自動整形" {
		return
	}
	
	//fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)

	var destination *string = nil
	var (
		output_S = ""
		output_O = ""
		output_A = ""
		output_P = ""
		output_C = ""
		output_I = ""
	)
	
	for _, v := range strings.Split(m.Content, "\n") {
		if strings.Contains(v, "主観的情報") {
			// start S
			output_S = "S:"
			destination = &output_S
		} else if strings.Contains(v, "客観的情報") {
			// start O
			output_O = "O:"
			destination = &output_O
		} else if checkAssesment(v) {
			// start A
			output_A = "A:"
			destination = &output_A
		} else if checkPlan(v) {
			// start P
			output_P = "P:"
			destination = &output_P
		} else if checkCarePlan(v) {
			// start C
			output_C = "C:"
			destination = &output_C
		} else if v == "備考" {
			// start info
			output_I = "備考:"
			destination = &output_I
		} else if destination != nil {
			*destination += v + "\n"
		}
	}
	/*
	fmt.Printf(output_S)
	fmt.Printf(output_O)
	fmt.Printf(output_A)
	fmt.Printf(output_P)
	fmt.Printf(output_C)
	fmt.Printf(output_I)
	*/

	sendReply(s, m.ChannelID, output_S, m.Reference())
	sendReply(s, m.ChannelID, output_O, m.Reference())
	sendReply(s, m.ChannelID, output_A, m.Reference())
	sendReply(s, m.ChannelID, output_P, m.Reference())
	sendReply(s, m.ChannelID, output_C, m.Reference())
	sendReply(s, m.ChannelID, output_I, m.Reference())
}

func checkAssesment(s string) bool {
/*	var s1, s2 string
	
	n, e := fmt.Sscanf(s, "%s(%s)", &s1, &s2)
	if(e != nil) {
		log.Printf("error: %s", e.Error())
		return false
	}
	if(n < 1) return false;
*/
	if strings.Contains(s, "評価") {
		if strings.Contains(s, "(A)") || strings.Contains(s, "（A）") ||
			 strings.Contains(s, "(評価)") || strings.Contains(s, "（評価）") {
			return true
		}
	}
	return false
}

func checkPlan(s string) bool {
	if strings.Contains(s, "計画") {
		if strings.Contains(s, "(P)") || strings.Contains(s, "（P）") ||
			 strings.Contains(s, "(計画)") || strings.Contains(s, "（計画）") {
			return true
		}
	}
	return false
}

func checkCarePlan(s string) bool {
	if strings.Contains(s, "ケア") {
		if strings.Contains(s, "(C)") || strings.Contains(s, "（C）") ||
			 strings.Contains(s, "(ケア)") || strings.Contains(s, "（ケア）") {
			return true
		}
	}
	return false
}

func sendMessage(s *discordgo.Session, channelID string, msg string) {
	 _, err := s.ChannelMessageSend(channelID, msg)
	 log.Println(">>> " + msg)
	 if err != nil {
		log.Println("Error sending message: ", err)
	 }
}

func sendReply(s *discordgo.Session, channelID string, msg string, reference *discordgo.MessageReference) {
	 _, err := s.ChannelMessageSendReply(channelID, msg, reference)
	 if err != nil {
		log.Println("Error sending message: ", err)
	 }
}

func listenPort() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	r.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Hello, %s!", name),
		})
	})

	r.Run(fmt.Sprintf(":%s", port))

}