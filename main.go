package main

func main() {
	// Get config
	cfg := new(Config)
	cfg.init()
	s := new(Server)
	s.init(cfg)
	//Scoop module
	ScoopModule := new(mod_scoop)
	ScoopModule.init(cfg)
	ScoopModule.init_route(s.api.Group("scoop"))
	SettingsModule := new(mod_settings)
	SettingsModule.init(cfg)
	SettingsModule.init_route(s.api.Group("settings"))
	s.start()
}
