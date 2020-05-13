package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/lancer-kit/uwe/v2"
	"github.com/lancer-kit/uwe/v2/presets/api"
	"github.com/sheb-gregor/polly-demo/mq"
	"github.com/sheb-gregor/polly-demo/server"
	"gopkg.in/yaml.v2"
)

type Config struct {
	API api.Config `yaml:"api"`
}

func main() {
	cfg := getConfiguration()

	chief := uwe.NewChief()
	// will add worker into the pool

	// will enable recover of internal panics
	chief.UseDefaultRecover()
	// pass handler for internal events like errors, panics, warning, etc.
	// you can log it with you favorite logger (ex Logrus, Zap, etc)
	chief.SetEventHandler(chiefEventHandler())

	broker := mq.NewBroker()
	chief.AddWorker("broker-server", api.NewServer(cfg.API, server.GetServer(broker)))

	// init all registered workers and run it all
	chief.Run()
}

func getConfiguration() (cfg Config) {
	path := flag.String("config", "config.yaml", "path to file with service configuration")
	flag.Parse()

	raw, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatal("FATAL: unable to read configuration; ", err.Error())
	}

	err = yaml.Unmarshal(raw, &cfg)
	if err != nil {
		log.Fatal("FATAL: unable to unmarshal configuration; ", err.Error())
	}
	return cfg
}

func chiefEventHandler() func(event uwe.Event) {
	return func(event uwe.Event) {
		var level string
		switch event.Level {
		case uwe.LvlFatal, uwe.LvlError:
			level = "ERROR"
		case uwe.LvlInfo:
			level = "INFO"
		default:
			level = "WARN"
		}
		log.Println(fmt.Sprintf("%s: %s %+v", level, event.Message, event.Fields))
	}
}
