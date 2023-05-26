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
	ScoopGroup := s.api.Group("scoop")
	ScoopModule.init_route(ScoopGroup)
	s.start()
}
