package main

func main() {
	// Get config
	cfg := new(Config)
	cfg.init()
	s := new(Server)
	s.init(cfg)
}
