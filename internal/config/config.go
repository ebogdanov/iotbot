package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	filePath string

	Tuya     *Tuya     `yaml:"tuya,omitempty"`
	Ewelink  *Ewelink  `yaml:"ewelink,omitempty"`
	Telegram *Telegram `yaml:"telegram"`
	Acl      *Acl      `yaml:"acl"`
	Stickers []string  `yaml:"stickers,omitempty"`
	Qr       *Qr       `yaml:"qr"`
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
	Users   map[string]string   `yaml:"users,omitempty"`
	Groups  map[string][]string `yaml:"groups,omitempty"`
	Invites []string            `yaml:"invites,omitempty"`
	Actions struct {
		Allow []string            `yaml:"allow"`
		Only  map[string][]string `yaml:"only,omitempty"`
	} `yaml:"actions"`
}

type Qr struct {
	Enable bool   `yaml:"enable,omitempty"`
	Cmd    string `yaml:"cmd"`
	Codes  []Code `yaml:"codes,omitempty"`
}

type Code struct {
	User  string `yaml:"user"`
	Title string `yaml:"title"`
	Code  string `yaml:"code"`
	Times int    `yaml:"times"`
}

// New returns a new decoded Config struct
func New(configPath string) (*Config, error) {
	var config *Config
	file, err := os.Open(configPath)
	if err == nil {
		defer func() { _ = file.Close() }()

		config = &Config{
			filePath: configPath,
		}

		err = yaml.NewDecoder(file).Decode(&config)
	}

	return config, err
}

func (c *Config) SaveFile() error {
	data, err := yaml.Marshal(c)
	if err == nil {
		err = os.WriteFile(c.filePath, data, 0644)
	}

	return err
}
