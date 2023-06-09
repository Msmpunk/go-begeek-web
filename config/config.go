package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Configuration struct {
	Db struct {
		Ip          string `json:"DbIp"`
		Port        string `json:"DbPort"`
		Name        string `json:"DbName"`
		User        string `json:"DbUser"`
		Password    string `json:"DbPassword"`
		AutoMigrate bool   `json:"DbAutoMigrate"`
	}

	Listen struct {
		Ip   string `json:"HttpIp"`
		Port string `json:"HttpPort"`
	}

	Template struct {
		BaseName string `json:"BaseTemplateName"`
	}
}

func LoadConfig() Configuration {
	c := flag.String("c", "env.json", "Path to the json configuration file")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("Unable to open JSON config file: ", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Unable to close JSON config file: ", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	Config := Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("Unable to decode JSON config file: ", err)
	}

	return Config
}
