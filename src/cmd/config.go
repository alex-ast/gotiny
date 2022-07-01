package main

import (
	"os"
	"path/filepath"

	"github.com/alex-ast/gotiny/cache"
	"github.com/alex-ast/gotiny/db"

	"github.com/spf13/viper"
)

type Config struct {
	CacheCfg   cache.CacheCfg     `mapstructure:"cache"`
	DbCfg      db.DbCfg           `mapstructure:"db"`
}

func LoadConfig(file string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Config search: current dir, dir with the executable, "../../conf" - for local dev build
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	currDir, _ := os.Getwd()

	viper.AddConfigPath(currDir)
	viper.AddConfigPath(exeDir)
	// For the development stage, lookup in source directory
	viper.AddConfigPath("../../conf")

	cfg := Config{}
	err := viper.ReadInConfig()
	if err == nil {
		err = viper.Unmarshal(&cfg)
	}
	return cfg, err
}
