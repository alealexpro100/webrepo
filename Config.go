package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	Port          string
	ClientPath    string
	AdminUser     string
	AdminPassword string
	ScoopEnabled  bool
	ScoopDir      string
	ConfigFile    *ini.File
}

func (cfg *Config) init() {
	var err error
	cfg.ConfigFile, err = ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	cfg.get()
}

func (cfg *Config) get() error {
	err := error(nil)
	cfg.Port = cfg.ConfigFile.Section("main").Key("web_port").String()
	cfg.ClientPath = cfg.ConfigFile.Section("main").Key("client_path").String()
	cfg.AdminUser = cfg.ConfigFile.Section("main").Key("user").String()
	cfg.AdminPassword = cfg.ConfigFile.Section("main").Key("password").String()
	cfg.ScoopEnabled, err = cfg.ConfigFile.Section("scoop").Key("enabled").Bool()
	cfg.ScoopDir = cfg.ConfigFile.Section("scoop").Key("directory").String()
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) set_auth(user, pass_hash string) {
	cfg.ConfigFile.Section("main").Key("user").SetValue(user)
	cfg.ConfigFile.Section("main").Key("password").SetValue(pass_hash)
	cfg.AdminUser = user
	cfg.AdminPassword = pass_hash
	cfg.ConfigFile.SaveTo("config.ini")
}
