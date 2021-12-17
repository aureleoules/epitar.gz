package api

import (
	"github.com/aureleoules/epitar/config"
	"github.com/expectedsh/go-sonic/sonic"
)

func initSonic() {
	var err error
	searcher, err = sonic.NewSearch(config.Cfg.Index.Sonic.Host, config.Cfg.Index.Sonic.Port, config.Cfg.Index.Sonic.Password)
	if err != nil {
		panic(err)
	}
}
