package main

func commandBoard(cfg *checkersCfg, params ...string) error{
	cfg.displayBoard()
	return nil
}