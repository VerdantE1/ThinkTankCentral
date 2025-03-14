package utils

import (
	"ThinkTankCentral/global"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const configFile = "config.yaml"

// LoadYAML 从文件中读取 YAML 数据并返回字节数组
func LoadYAML() ([]byte, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil, err
	}
	return data, nil
}

// SaveYAML 将全局配置对象保存为 YAML 格式到文件
func SaveYAML(data []byte) error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		log.Printf("Error marshalling yaml: %v", err)
		return err
	}
	return os.WriteFile(configFile, byteData, 0644)
}
