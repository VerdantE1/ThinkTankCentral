package config

type Config struct {
	Captcha Captcha `json:"captcha" yaml:"captcha"`
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`
	Zap     Zap     `json:"zap" yaml:"zap"`
	System  System  `json:"system" yaml:"system"`
	ES      ES      `json:"es" yaml:"es"`
	Email   Email   `json:"email" yaml:"email"`
	Redis   Redis   `json:"redis" yaml:"redis"`
	Jwt     Jwt     `json:"jwt" yaml:"jwt"`
}
