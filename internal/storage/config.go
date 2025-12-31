package storage

import "gorm.io/gorm"

type Config struct {
	DB   *gorm.DB `yaml:"-" json:"-"`
	Port int      `yaml:"port" json:"port"`

	EndPoint        string `yaml:"endpoint" json:"endpoint"`
	Region          string `yaml:"region" json:"region"`
	AccessKeyID     string `yaml:"accessKeyID" json:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret" json:"accessKeySecret"`
	Bucket          string `yaml:"bucket" json:"bucket"`
	ForcePathStyle  bool   `yaml:"forcePathStyle" json:"forcePathStyle"`
}

func NewConfig() *Config {
	return &Config{
		Port:            8082,
		EndPoint:        "",
		Region:          "us-east-1",
		AccessKeyID:     "",
		AccessKeySecret: "",
		Bucket:          "",
		ForcePathStyle:  false,
	}
}
