package config

type Config struct {
	Captcha Captcha `json:"captcha" yaml:"captcha"`
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`
	Zap     Zap     `json:"zap" yaml:"zap"`
	System  System  `json:"system" yaml:"system"`
}
