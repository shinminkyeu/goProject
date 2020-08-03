package config

import (
	"os"
	"path"

	"github.com/klaytn/klaytn/common"
	"github.com/naoina/toml"
)

//Config Config
type Config struct {
	Keystore struct {
		Path     string
		Owner    common.Address
		Feepayer common.Address
	}

	ContractDB struct {
		URL        string
		DB         string
		Collection string
	}
}

//NewConfig NewConfig
func NewConfig(file string) (*Config, error) {
	c := new(Config)
	if file, err := os.Open(file); err != nil {
		return nil, err
	} else {
		defer file.Close()
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			return nil, err
		} else if err := c.sanitize(); err != nil {
			return nil, err
		} else {
			return c, nil
		}
	}
}

func (p *Config) sanitize() error {
	if p.Keystore.Path[0] == byte('~') {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		p.Keystore.Path = path.Join(home, p.Keystore.Path[1:])
	}
	return nil
}
