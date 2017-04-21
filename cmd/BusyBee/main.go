package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"io/ioutil"

	"log"

	"github.com/dixonwille/busybee"
	_ "github.com/dixonwille/busybee/exchange"
	_ "github.com/dixonwille/busybee/hipchat"
	"github.com/dixonwille/busybee/util"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

func main() {
	cfg, err := parseConfig("busybee.yml")
	if err != nil {
		log.Fatalln(err)
	}
	eventService, statusService, err := createServices(cfg.Plugins)
	if err != nil {
		log.Fatalln(err)
	}
	user := busybee.NewUser(cfg.EventUID, cfg.StatusUID, eventService, statusService)
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

func askQuestions(conf interface{}, askAll bool) (bool, error) {
	changed := false
	confV := reflect.ValueOf(conf)
	if confV.Kind() != reflect.Ptr {
		return changed, errors.New("Configuration must be a pointer to a struct")
	}
	v := confV.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)
		//Make sure it is not already populated
		lengther := fv.Kind() == reflect.Map || fv.Kind() == reflect.Slice || fv.Kind() == reflect.Chan || fv.Kind() == reflect.Array || fv.Kind() == reflect.String
		var empty bool
		if lengther {
			empty = fv.Len() == 0
		} else {
			empty = fv.Interface() == reflect.Zero(fv.Type()).Interface()
		}
		if (empty || askAll) && ft.Tag.Get("quest") != "" {
			args := strings.Split(ft.Tag.Get("quest"), ",")
			if contains("encrypt", args) && fv.Kind() != reflect.String {
				return changed, fmt.Errorf("Struct: %s Type: %s must be of type string if you want it encrypted", t.Name(), ft.Name)
			}
			err := askAndUpdate(fv, args[0], contains("encrypt", args), contains("pass", args))
			if err != nil {
				return changed, err
			}
			changed = true
		}
	}
	return changed, nil
}

func contains(find string, in []string) bool {
	for _, i := range in {
		if i == find {
			return true
		}
	}
	return false
}

func askAndUpdate(v reflect.Value, question string, encrypted, password bool) error {
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
		if res == "" {
			continue
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

func parseConfig(cfg string) (*busybee.MainConfig, error) {
	conf := new(busybee.MainConfig)
	var fileBytes []byte
	if file, err := os.Open(cfg); err == nil {
		defer file.Close()
		fileBytes, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	} else {
		fileBytes = make([]byte, 0)
	}
	err := yaml.Unmarshal(fileBytes, conf)
	if err != nil {
		return nil, err
	}
	mainChanged, err := askQuestions(conf, false)
	if err != nil {
		return nil, err
	}
	var eventChanged bool
	var statusChanged bool
	var eventServiceName string
	var statusServiceName string
	var eventConf interface{}
	var statusConf interface{}
ParsePlugins:
	for name, c := range conf.Plugins {
		switch c.Type {
		case busybee.ServiceTypeEvent:
			if eventServiceName != "" {
				break
			}
			eventServiceName = name
			eventService, err := busybee.GetEventService(name)
			if err != nil {
				return nil, err
			}
			eventConf = eventService.CreateConfig()
			eConfMarsh, err := yaml.Marshal(c.Config)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(eConfMarsh, eventConf)
			if err != nil {
				return nil, err
			}
			eventChanged, err = askQuestions(eventConf, false)
			if err != nil {
				return nil, err
			}
		case busybee.ServiceTypeStatus:
			if statusServiceName != "" {
				break
			}
			statusServiceName = name
			statusService, err := busybee.GetStatusService(name)
			if err != nil {
				return nil, err
			}
			statusConf = statusService.CreateConfig()
			sConfMarsh, err := yaml.Marshal(c.Config)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(sConfMarsh, statusConf)
			if err != nil {
				return nil, err
			}
			statusChanged, err = askQuestions(statusConf, false)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("Do not know what type %d is", c.Type)
		}
	}
	if eventServiceName == "" {
		var name string
		err = askAndUpdate(reflect.ValueOf(&name).Elem(), "Which Event Service plugin would you like to use?", false, false)
		if err != nil {
			return nil, err
		}
		if conf.Plugins == nil {
			conf.Plugins = make(map[string]busybee.PluginConfig)
		}
		conf.Plugins[name] = busybee.PluginConfig{
			Type: busybee.ServiceTypeEvent,
		}
		goto ParsePlugins
	}
	if statusServiceName == "" {
		var name string
		err = askAndUpdate(reflect.ValueOf(&name).Elem(), "Which Status Service plugin would you like to use?", false, false)
		if err != nil {
			return nil, err
		}
		if conf.Plugins == nil {
			conf.Plugins = make(map[string]busybee.PluginConfig)
		}
		conf.Plugins[name] = busybee.PluginConfig{
			Type: busybee.ServiceTypeStatus,
		}
		goto ParsePlugins
	}
	returnCfg := new(busybee.MainConfig)
	returnCfg.EventUID = conf.EventUID
	returnCfg.StatusUID = conf.StatusUID
	returnCfg.Plugins = make(map[string]busybee.PluginConfig, 2)
	returnCfg.Plugins[eventServiceName] = busybee.PluginConfig{
		Type:   busybee.ServiceTypeEvent,
		Config: eventConf,
	}
	returnCfg.Plugins[statusServiceName] = busybee.PluginConfig{
		Type:   busybee.ServiceTypeStatus,
		Config: statusConf,
	}
	if mainChanged || eventChanged || statusChanged {
		newFileBytes, err := yaml.Marshal(returnCfg)
		if err != nil {
			return nil, err
		}
		if err = ioutil.WriteFile(cfg, newFileBytes, 0600); err != nil {
			return nil, err
		}
	}
	return returnCfg, nil
}

func createServices(plugins map[string]busybee.PluginConfig) (busybee.InEventer, busybee.UpdateStatuser, error) {
	var event busybee.InEventer
	var status busybee.UpdateStatuser
	for name, plugin := range plugins {
		switch plugin.Type {
		case busybee.ServiceTypeEvent:
			eventService, err := busybee.GetEventService(name)
			if err != nil {
				return nil, nil, err
			}
			event, err = eventService.Create(plugin.Config)
			if err != nil {
				return nil, nil, err
			}
		case busybee.ServiceTypeStatus:
			statusService, err := busybee.GetStatusService(name)
			if err != nil {
				return nil, nil, err
			}
			status, err = statusService.Create(plugin.Config)
			if err != nil {
				return nil, nil, err
			}
		default:
			return nil, nil, fmt.Errorf("Do not know how to create an event for type %d", plugin.Type)
		}
	}
	return event, status, nil
}
