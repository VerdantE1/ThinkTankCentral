package config

type Config struct {
	Captcha Captcha `json:"captcha" yaml:"captcha"`
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`
}
