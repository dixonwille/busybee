package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"flag"

	"github.com/dixonwille/busybee"
	_ "github.com/dixonwille/busybee/exchange"
	_ "github.com/dixonwille/busybee/hipchat"
	"github.com/dixonwille/busybee/util"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	cfgFile string
)

func init() {
	flag.StringVar(&cfgFile, "cfg", "busybee.yaml", "Specify where to find the configuration file.")
	flag.Parse()
}

func main() {
	event, err := createInEventer("exchange")
	if err != nil {
		log.Fatalln(err)
	}
	status, err := createUpdateStatuser("hipchat")
	if err != nil {
		log.Fatalln(err)
	}
	user := busybee.NewUser("", "", status, event)
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
	if err = askQuestions(serviceConf, false); err != nil {
		return nil, err
	}
	return eventService.Create(serviceConf)
}

func createUpdateStatuser(name string) (busybee.UpdateStatuser, error) {
	statusService, err := busybee.GetStatusService(name)
	if err != nil {
		return nil, err
	}
	serviceConf := statusService.CreateConfig()
	if err = askQuestions(serviceConf, false); err != nil {
		return nil, err
	}
	return statusService.Create(serviceConf)
}

func askQuestions(conf interface{}, askAll bool) error {
	confV := reflect.ValueOf(conf)
	if confV.Kind() != reflect.Ptr {
		return errors.New("Configuration must be a pointer to a struct")
	}
	v := confV.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)
		//Make sure it is not already populated
		if (fv.Interface() == reflect.Zero(fv.Type()).Interface() || askAll) && ft.Tag.Get("quest") != "" {
			args := strings.Split(ft.Tag.Get("quest"), ",")
			if contains("encrypt", args) && fv.Kind() != reflect.String {
				return fmt.Errorf("Struct: %s Type: %s must be of type string if you want it encrypted", t.Name(), ft.Name)
			}
			err := askAndUpdate(fv, args[0], contains("required", args), contains("encrypt", args), contains("pass", args))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func contains(find string, in []string) bool {
	for _, i := range in {
		if i == find {
			return true
		}
	}
	return false
}

func askAndUpdate(v reflect.Value, question string, required, encrypted, password bool) error {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)
	term := terminal.NewTerminal(os.Stdout, "")
question:
	for {
		term.Write([]byte(util.CleanQuest(question)))
		var res string
		if password {
			res, err = term.ReadPassword("")
		} else {
			res, err = term.ReadLine()
		}
		if err != nil {
			return err
		}
		res = strings.Replace(res, "\r", "", -1)
		res = strings.Replace(res, "\n", "", -1)
		res = strings.Trim(res, " ")
		if res == "" && required {
			continue
		}
		if res == "" {
			break
		}
		switch v.Kind() {
		case reflect.String:
			v.SetString(res) //Do something on encryption
			break question
		default:
			fmt.Printf("Not sure how to convert string to type %s\n", v.Kind().String())
		}
	}
	return nil
}
