package main

import (
	"fmt"
	"os"
	"wild_bluebell/controller"
	"wild_bluebell/dao/mysql"
	"wild_bluebell/dao/redis"
	"wild_bluebell/logger"
	"wild_bluebell/pkg/snowflake"
	"wild_bluebell/router"
	"wild_bluebell/setting"
)

func main() {
	// 命令行参数检查，确保用户启动程序时提供了配置文件路径
	// os.Args[0]:程序自身路径，例如 ./wild_bluebell
	// os.Args[1]:第一个命令行参数，例如 config.yaml
	// os.Args[2]: 第二个命令行参数
	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: wild_bluebell config.yaml")
		return
	}

	// 加载配置文件
	if err := setting.Init(os.Args[1]); err != nil {
		fmt.Printf("load config fialed, err:%v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 初始化MySQL数据库连接
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() //程序退出时关闭数据库连接

	// 初始化Redis数据库连接
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 初始化雪花算法，生成ID
	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator translate failed, err:%v\n", err)
		return
	}

	// 注册路由
	r := router.SetupRouter(setting.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
