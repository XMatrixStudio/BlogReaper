package app

import "github.com/XMatrixStudio/Violet.SDK.Go"

// Config 配置文件
type Config struct {
	Violet violetSdk.Config `yaml:"Violet"` // Violet配置
	Port   int              `yaml:"Port"`   // 端口
}
