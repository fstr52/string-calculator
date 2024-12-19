package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Addr string `json:"addr"`
}

func GetConfig(filePath string) {
	var config Config
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}

		defer file.Close()

		config = Config{Addr: "8080"}

		jsonStockConfig, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}

		_, err = file.Write(jsonStockConfig)
		if err != nil {
			panic(err)
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			panic(err)
		}

	} else {
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		if len(jsonData) == 0 {
			panic(fmt.Errorf("config file is empty"))
		}

		err = json.Unmarshal(jsonData, &config)
		if err != nil {
			panic(fmt.Errorf("config unmarshaling failed, error: %s", err))
		}
	}

	err := os.Setenv("PORT", config.Addr)
	if err != nil {
		panic(err)
	}
}
