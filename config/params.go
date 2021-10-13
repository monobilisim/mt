package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"mt/log"
	"mt/minio"
	"os"
)

type Params struct {
	Minio map[string]minio.Params
	Log log.Params
}

func NewParams(config string) (p *Params, err error) {
	if _, err = os.Stat(config); os.IsNotExist(err) {
		fmt.Printf("Configuration file: %s does not exist, %v\n", config, err)
		return nil, err
	}

	viper.SetConfigFile(config)
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Unable to read config file, %v\n", err)
		return nil, err
	}

	err = viper.Unmarshal(&p)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
		return nil, err
	}

	return
}

func (p *Params) Server(server string) (minio.Params, error) {
	if s, ok := p.Minio[server]; ok {
		return s, nil
	}
	return minio.Params{}, errors.New("no server found with key")
}
