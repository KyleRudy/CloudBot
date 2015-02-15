// cloudbot project main.go
package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
)

func performHostCheck(hostname string, tor *TorCache, ircObj *irc.Connection, e *irc.Event) {
	result, err := tor.Check(hostname)
	if err != nil {
		ircObj.Privmsg(e.Arguments[0], "Error encountered: "+err.Error())
	} else if result != "" {
		ircObj.Privmsg(e.Arguments[0], hostname+" was found on the list of TOR exit nodes in "+result)
	} else {
		ircObj.Privmsg(e.Arguments[0], hostname+" looks clean.")
	}
}

func main() {
	fmt.Println("Initializing CloudBot.")
	config, err := RetrieveConfig()
	if err != nil {
		fmt.Println("Error encountered while loading conf.json:")
		fmt.Println(err.Error())
		return
	}
	tor := CreateTorCache(strings.Split(config.LookupLocations, ";"))
	ircObj := irc.IRC(config.Nick, config.User)
	ircObj.Connect(config.Server)
	ircObj.VerboseCallbackHandler = true
	ircObj.Debug = true

	ircObj.AddCallback("INVITE", func(e *irc.Event) {
		ircObj.Join(e.Arguments[0])
	})

	ircObj.AddCallback("PRIVMSG", func(e *irc.Event) {
		// ircObj.Privmsg(e.Arguments[0], e.Message()+" to you too, "+e.Nick)
		msg := e.Message()
		terms := strings.Split(msg, " ")
		switch terms[0] {
		case "!check":
			for _, element := range terms[1:] {
				if len(element) > 0 {
					performHostCheck(element, tor, ircObj, e)
				}
			}
		case "!leave":
		case "!bye":
		case "!quit":
			ircObj.Part(e.Arguments[0])
		case "!join":
		case "!invite":
			for _, element := range terms[1:] {
				if len(element) > 0 {
					ircObj.Join(element)
				}
			}
		case "!num":
			ircObj.Privmsgf(e.Arguments[0], "Currently caching %d IPs!", tor.GetNumberOfIPs())
		}
	})

	ircObj.AddCallback("001", func(e *irc.Event) {
		ircObj.Join("#cloudbottest")
		ircObj.SendRaw("oper " + config.Nick + " " + config.OperPwd)
	})

	ircObj.Loop()

	fmt.Println("Hello again")
}
