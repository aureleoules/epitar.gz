package api

import (
	"time"

	"github.com/aureleoules/epitar/config"
	"github.com/expectedsh/go-sonic/sonic"
	"github.com/silenceper/pool"
)

var sonicPool pool.Pool

func initSonic() {
	var err error
	sonicPool, err = pool.NewChannelPool(&pool.Config{
		InitialCap: 5,
		MaxCap:     30,
		MaxIdle:    20,
		Factory: func() (interface{}, error) {
			return sonic.NewSearch(config.Cfg.Index.Sonic.Host, config.Cfg.Index.Sonic.Port, config.Cfg.Index.Sonic.Password)
		},
		Close: func(v interface{}) error {
			return v.(sonic.Searchable).Quit()
		},
		Ping: func(v interface{}) error {
			return v.(sonic.Searchable).Ping()
		},
		IdleTimeout: time.Second * 10,
	})

	if err != nil {
		panic(err)
	}
}
