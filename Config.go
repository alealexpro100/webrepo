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
	ScoopInitRepo string
	AptlyEnabled  bool
	AptlyLocal    bool
	AptlyURL      string
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
	// For scoop
	cfg.ScoopEnabled, err = cfg.ConfigFile.Section("scoop").Key("enabled").Bool()
	if err != nil {
		return err
	}
	cfg.ScoopDir = cfg.ConfigFile.Section("scoop").Key("directory").String()
	cfg.ScoopInitRepo = cfg.ConfigFile.Section("scoop").Key("init_repo").String()
	// For aptly
	cfg.AptlyEnabled, err = cfg.ConfigFile.Section("aptly").Key("enabled").Bool()
	if err != nil {
		return err
	}
	cfg.AptlyLocal, err = cfg.ConfigFile.Section("aptly").Key("local").Bool()
	if err != nil {
		return err
	}
	cfg.AptlyURL = cfg.ConfigFile.Section("aptly").Key("url").String()
	return nil
}

func (cfg *Config) set_auth(user, pass_hash string) {
	cfg.ConfigFile.Section("main").Key("user").SetValue(user)
	cfg.ConfigFile.Section("main").Key("password").SetValue(pass_hash)
	cfg.AdminUser = user
	cfg.AdminPassword = pass_hash
	cfg.ConfigFile.SaveTo("config.ini")
}
