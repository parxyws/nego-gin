package main

import (
	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/pkg/database/psql"
)

func main() {
	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	_, err = psql.NewDB(cfg)
	if err != nil {
		panic(err)
	}

}
