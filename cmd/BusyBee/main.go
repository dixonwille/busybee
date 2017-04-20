package main

import (
	"flag"
	"fmt"

	"log"

	"os"

	"strings"

	"github.com/dixonwille/busybee"
	"github.com/dixonwille/busybee/exchange"
	"github.com/dixonwille/busybee/hipchat"
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
	event, err := createInEventer("exchange")
	if err != nil {
		log.Fatalln(err)
	}
	status, err := createUpdateStatuser("hipchat")
	user := busybee.NewUser(exUID, cleanMention(hcUID), status, event)
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
func createInEventer(name string) (busybee.InEventer, error) {
	eventService, err := busybee.GetEventService(name)
	if err != nil {
		return nil, err
	}
	serviceConf := eventService.CreateConfig()
	//TODO make sure that this is where I load for configuration file.
	//Should not need to convert the struct as I do not care.
	exchangeConf, ok := serviceConf.(*exchange.Conf)
	if !ok {
		return nil, fmt.Errorf("Could not convert the configuration struct to a %s configuration struct", name)
	}
	exchangeConf.Host = cleanHost(exHost)
	exchangeConf.Pass = exPass
	exchangeConf.User = exUser
	return eventService.Create(exchangeConf)
}

func createUpdateStatuser(name string) (busybee.UpdateStatuser, error) {
	statusService, err := busybee.GetStatusService(name)
	if err != nil {
		return nil, err
	}
	serviceConf := statusService.CreateConfig()
	//TODO make sure that this is where I load for configuration file.
	//Should not need to convert the struct as I do not care.
	hipchatConf, ok := serviceConf.(*hipchat.Conf)
	if !ok {
		return nil, fmt.Errorf("Could not convert the configuration struct to a %s configuration struct", name)
	}
	hipchatConf.Host = cleanHost(hcHost)
	hipchatConf.Token = hcToken
	return statusService.Create(hipchatConf)
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
