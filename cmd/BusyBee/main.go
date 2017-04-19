package main

import (
	"flag"

	"log"

	"os"

	"strings"

	"github.com/dixonwille/busybee"
	_ "github.com/dixonwille/busybee/exchange"
	_ "github.com/dixonwille/busybee/hipchat"
)

var (
	exUID   string
	exHost  string
	exUser  string
	exPass  string
	hcUID   string
	hcHost  string
	hcToken string
)

func init() {
	flag.StringVar(&exUID, "exUID", "", "The UID of the user for Exchange (usually an email address)")
	flag.StringVar(&exHost, "exHost", "", "The host of the Exchange server")
	flag.StringVar(&exUser, "exUser", "", "The user to sign in as when checking the Exchange server")
	flag.StringVar(&exPass, "exPass", "", "The password of the user to sign in as when checkint the Exchange server")
	flag.StringVar(&hcUID, "hcUID", "", "The UID of the user for Hipchat (usually @mention name)")
	flag.StringVar(&hcHost, "hcHost", "", "The host of the hipchat server (uaually team.hipchat.com)")
	flag.StringVar(&hcToken, "hcToken", "", "The token to use to validate calls")
	flag.Parse()
	wasErr := false
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			log.Printf("%s must not be empty", f.Name)
			wasErr = true
		}
	})
	if wasErr {
		os.Exit(1)
	}
}

func main() {
	calSerConf := make(map[string]string)
	calSerConf["service"] = "exchange"
	calSerConf["host"] = cleanHost(exHost)
	calSerConf["user"] = exUser
	calSerConf["pass"] = exPass
	staSerConf := make(map[string]string)
	staSerConf["service"] = "hipchat"
	staSerConf["host"] = cleanHost(hcHost)
	staSerConf["token"] = hcToken
	calendar, err := busybee.CreateCalendarService(calSerConf)
	if err != nil {
		log.Fatalln(err)
	}
	status, err := busybee.CreateStatusService(staSerConf)
	if err != nil {
		log.Fatalln(err)
	}
	user := busybee.NewUser(exUID, cleanMention(hcUID), status, calendar)
	inEvent, err := user.InEvent()
	if err != nil {
		log.Fatalln(err)
	}
	curStatus := busybee.StatusUnknown
	if inEvent {
		curStatus = busybee.StatusBusy
	} else {
		curStatus = busybee.StatusAvailable
	}
	err = user.UpdateStatus(curStatus)
	if err != nil {
		log.Fatalln(err)
	}
}

func cleanHost(host string) string {
	if strings.Index(host, "http://") == 0 || strings.Index(host, "https://") == 0 {
		return host
	}
	return "https://" + host
}

func cleanMention(mention string) string {
	mentionRunes := []rune(mention)
	if mention[0] != '@' {
		mentionRunes = append([]rune{rune('@')}, mentionRunes...)
	}
	return string(mentionRunes)
}
