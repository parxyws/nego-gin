package test

import (
	"os"

	"github.com/parxyws/nego-gin/config"
)

func InitEnv() (*config.Config, error) {
	paths := []string{"config.yaml", "../config.yaml", "../../config.yaml", "../../../config.yaml"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			_ = os.Setenv("CONFIG_PATH", p)
			break
		}
	}
	cfg, err := config.NewAppConfig()
	return cfg, err
}
