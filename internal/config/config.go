package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Db       *Db       `yaml:"db"`
	Tuya     *Tuya     `yaml:"tuya,omitempty"`
	Ewelink  *Ewelink  `yaml:"ewelink,omitempty"`
	Telegram *Telegram `yaml:"telegram"`
	Acl      *Acl      `yaml:"acl"`
	Stickers []string  `yaml:"stickers,omitempty"`
	Qr       *Qr       `yaml:"qr"`
}

type Db struct {
	Driver     string `yaml:"driver"`
	Connection string `yaml:"connection"`
}

type Tuya struct {
	AccessId  string `yaml:"access_id"`
	AccessKey string `yaml:"access_key"`
	UserId    string `yaml:"user_id"`
}

type Ewelink struct {
	Region   string `yaml:"region"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Telegram struct {
	TokenBot string `yaml:"token_bot"`
}

type Acl struct {
	Actions struct {
		Allow []string            `yaml:"allow"`
		Only  map[string][]string `yaml:"only,omitempty"`
	} `yaml:"actions"`
}

type Qr struct {
	Enable     bool   `yaml:"enable,omitempty"`
	Cmd        string `yaml:"cmd"`
	AllowCodes bool   `yaml:"allow_codes,omitempty"`
}

// New returns a new decoded Config struct
func New(configPath string) (*Config, error) {
	var config *Config
	file, err := os.Open(configPath)
	if err == nil {
		defer func() { _ = file.Close() }()

		config = &Config{}

		err = yaml.NewDecoder(file).Decode(&config)
	}

	return config, err
}
