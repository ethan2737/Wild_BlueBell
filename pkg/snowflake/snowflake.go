// Package snowflake 雪花算法：采用分布式生成全局唯一ID，支持大并发
// 该算法生成的ID具有以下特点：
// 1. 分布式环境下全局唯一
// 2. 趋势递增（适合数据库索引）
// 3. 高并发支持
// 4. 包含时间戳信息
package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

// node 是雪花算法的工作节点
// 每个节点在分布式环境中拥有唯一的机器ID
var node *sf.Node

// Init 初始化雪花算法节点
// 参数说明：
//   - startTime: 起始时间字符串，格式为 "2006-01-02"，用于计算ID的时间戳基准
//   - machineID: 机器ID（节点ID），取值范围 0~1023，每个分布式节点应有唯一的ID
//
// 返回值：
//   - err: 初始化过程中的错误，如时间解析失败或节点创建失败
//
// 使用示例：
//
//	if err := snowflake.Init("2024-01-01", 1); err != nil {
//	    log.Fatal(err)
//	}
//	id := snowflake.GenID()
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	// 解析起始时间，格式：yyyy-mm-dd
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	// 设置雪花算法的epoch（起始时间点），单位：毫秒
	sf.Epoch = st.UnixNano() / 1000000
	// 创建雪花算法节点，machineID 必须在 0~1023 范围内
	node, err = sf.NewNode(machineID)
	return
}

// GenID 生成一个全局唯一的分布式ID
//
// 返回值：
//   - int64: 生成的唯一ID
//
// 生成的ID结构（从高位到低位）：
//   - 41位时间戳
//   - 10位机器ID
//   - 12位序列号
//
// 注意事项：
//   - 必须先调用 Init 初始化后才能使用
//   - 每次调用 GenID 前必须确保已调用 Init
func GenID() int64 {
	return node.Generate().Int64()
}
