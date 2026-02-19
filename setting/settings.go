// Package setting 配置管理模块，使用Viper库从配置文件读取应用配置，并映射到结构体变量Conf
package setting

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	// mapstructure 是 Viper 专用的标签，用于将配置文件中的键名映射到 Go 结构体字段
	Name      string `mapstructure:"name"` // 配置文件中的 "name" → 映射到 Name
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	Port         int    `mapstructure:"port"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// Init 初始化函数，加载setting包后，自动执行 Init 初始化函数
func Init(filePath string) (err error) {
	// 指定配置文件路径（支持 YAML、JSON、TOML 等格式）
	viper.SetConfigFile(filePath)

	// 读取配置信息
	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 conf 结构体
	// Unmarshal 将配置解析到结构体（依赖 mapstructure 标签）
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	// 监听配置文件变化（热更新）
	viper.WatchConfig()
	//  回调函数：文件变化后重新 Unmarshl，实现配置热更新
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了……")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}

/*
工作流程图

  config.yaml
      ↓
  ReadInConfig()  →  Viper 内部存储
      ↓
  Unmarshal()     →  conf 结构体
      ↓
  运行时文件变化 → WatchConfig 触发回调 → 重新 Unmarshal → 配置更新

*/
