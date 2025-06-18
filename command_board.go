package main

func commandBoard(cfg *checkersCfg, params ...string) error{
	return cfg.displayBoard()
}