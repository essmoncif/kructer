package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Address          string `env:"ADDRESS" envDefault:":8080"`
	ConnectionString string `env:"CONNECTION_STRING,required"`
	Dialect          string `env:"DIALECT,required" envDefault:"postgres"`
	Secret           string `env:"SECRET,required"`
}

func New(files ...string) (*Configuration, error) {
	if err := godotenv.Load(files...); err != nil {
		panic("No env file could be found")
	}

	conf := Configuration{}
	if err := env.Parse(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
