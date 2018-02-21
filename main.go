package main

import (
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/motemen/go-pocket/api"
)

var configDir string
var k2pdfopt string

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configDir = filepath.Join(usr.HomeDir, ".config", "if-pocket-then-kindle")
	err = os.MkdirAll(configDir, 0777)
	if err != nil {
		panic(err)
	}

	k2pdfopt = "k2pdfopt"
}

func main() {
	bot := newBot()
	bot.run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	bot.shutdown()
}

type Bot struct {
	waitGroup sync.WaitGroup
	stop      chan struct{}
}

func newBot() *Bot {
	return &Bot{
		stop: make(chan struct{}),
	}
}

func (bot *Bot) run() {
	log.Println("Running...")
	consumerKey := getConsumerKey()
	accessToken, err := restoreAccessToken(consumerKey)
	if err != nil {
		panic(err)
	}

	client := api.NewClient(consumerKey, accessToken.AccessToken)

	mail, err := restoreMailSettings()
	if err != nil {
		panic(err)
	}

	since := time.Now().Unix()
	timeout := time.Duration(time.Second * 10)
	for {
		select {
		case <-bot.stop:
			return
		case <-time.After(timeout):
			if items, err := commandList(client, since); err != nil {
				log.Println(err)
			} else {
				bot.waitGroup.Add(1)
				go func() {
					defer bot.waitGroup.Done()
					bot.handle(items, mail)
				}()
			}
			since = time.Now().Unix()
		}
	}
}

func (bot *Bot) handle(items []api.Item, mail *Mail) {
	for _, item := range items {
		fp := mkdir(item)
		download(fp, item)
		converted := convert(fp, item)
		if err := sendToKindle(mail, converted); err != nil {
			log.Println(err)
		}
	}
}

func (bot *Bot) shutdown() {
	close(bot.stop)
	bot.waitGroup.Wait()
}
