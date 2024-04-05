package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"waba-alertmanager/notify"
	"waba-alertmanager/providers/acs"
)

var (
	listenAddress = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")
	configFile    = flag.String("config", "config.yaml", "The configuration file")
)


func main() {
	flag.Parse()

	if err := LoadConfig(*configFile); err != nil {
		fmt.Errorf("Error loading configuration: %s", err)
	}

	app := handlers{}

	http.HandleFunc("/alerts", app.Alert)

	if os.Getenv("PORT") != "" {
		*listenAddress = ":" + os.Getenv("PORT")
	}

	fmt.Printf("Listening on %s", *listenAddress)

	http.ListenAndServe(*listenAddress, nil)
}

func receiverConfByReceiver(name string) *ReceiverConf {
	for i := range config.Receivers {
		rc := &config.Receivers[i]
		if rc.Name == name {
			return rc
		}
	}
	return nil
}

func providerByName(name string) (notify.Provider, error) {
	switch name {
	case "acs":
		return acs.NewAcs(config.Providers.Acs), nil
	}

	return nil, fmt.Errorf("%s: Unknown provider", name)
}

func errorHandler(w http.ResponseWriter, status int, err error, provider string) {
	w.WriteHeader(status)

	data := struct {
		Error   bool
		Status  int
		Message string
	}{
		true,
		status,
		err.Error(),
	}
	// respond json
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("marshalling error: " + err.Error())
	}

	if _, err := w.Write(body); err != nil {
		fmt.Errorf("marshalling error: " + err.Error())
	}

	fmt.Println("error: " + string(body))
}
