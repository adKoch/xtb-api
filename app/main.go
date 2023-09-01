package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/adKoch/xtb-api/app/lib/config"
	"github.com/adKoch/xtb-api/app/lib/log"
	"github.com/gorilla/websocket"
) //import "os"

const (
	configPath      = "."
	configFile      = "properties"
	configExtension = "env"
)

func main() {
	// Load config
	config.LoadConfig(configPath, configFile, configExtension)

	// Init websocket
	ad := config.GetConfig("XTB_API_URL")
	log.Info(fmt.Sprintf("connecting to %s", ad))
	c, _, err := websocket.DefaultDialer.Dial(ad, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("dial: %s", err))
	}
	defer c.Close()

	// On read data
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Info(fmt.Sprintf("read error: %s", err))
				return
			}
			log.Info(fmt.Sprintf("recv: %s", message))
		}
	}()

	// Log in
	loginCmd := fmt.Sprintf(`{
		"command": "login",
		"arguments": {
			"userId": "%s",
			"password": "%s",
			"appName": "test"
		}
	}`, config.GetConfig("XTB_USERID"), config.GetConfig("XTB_PASSWORD"))

	err = c.WriteMessage(websocket.TextMessage, []byte(loginCmd))
	if err != nil {
		log.Info(fmt.Sprint("write: ", err))
		return
	} else {
		log.Info(fmt.Sprint("Logging in: ", loginCmd))
	}

	// Create time ticker
	prd, err := strconv.ParseInt(config.GetConfig("TICK_INTERVAL_SECONDS"), 10, 32)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not convert tick interval property to string error: %s", err))
	}
	ticker := time.NewTicker(time.Second * time.Duration(prd))
	defer ticker.Stop()

	msg := `{
		"command": "getChartLastRequest",
		"arguments": {
			"info": {
				"period": 1440,
				"start": 1,
				"symbol": "US100"
			}
		}
	}
	`
	// Send data every tick
	// interrupt when sent signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(msg)) //[]byte(t.String()))
			if err != nil {
				log.Info(fmt.Sprintf("write: %s", err))
				return
			} else {
				log.Info(fmt.Sprintf("write: %s", msg))
			}
		case <-interrupt:
			log.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Info(fmt.Sprintf("write close: %s", err))
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
