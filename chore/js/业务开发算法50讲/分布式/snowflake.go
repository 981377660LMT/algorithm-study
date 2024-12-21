package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net"
	"sync"
	"time"
)

const (
	UNUSED_BITS   = 1
	EPOCH_BITS    = 41
	NODE_ID_BITS  = 10
	SEQUENCE_BITS = 12

	// Custom Epoch (January 1, 2015 Midnight UTC = 2015-01-01T00:00:00Z)
	DEFAULT_CUSTOM_EPOCH = 1420070400000
)

var (
	maxNodeId   = (1 << NODE_ID_BITS) - 1
	maxSequence = (1 << SEQUENCE_BITS) - 1
)

// Snowflake 结构体类似于 Java 类，保存状态信息
type Snowflake struct {
	mu            sync.Mutex
	nodeId        int64
	customEpoch   int64
	lastTimestamp int64
	sequence      int64
}

// 使用指定的 nodeId 和 customEpoch 创建 Snowflake 实例
func NewSnowflake(nodeId, customEpoch int64) *Snowflake {
	if nodeId < 0 || nodeId > int64(maxNodeId) {
		panic(fmt.Sprintf("NodeId must be between %d and %d", 0, maxNodeId))
	}
	return &Snowflake{
		nodeId:        nodeId,
		customEpoch:   customEpoch,
		lastTimestamp: -1,
		sequence:      0,
	}
}

// 使用指定的 nodeId 和默认 epoch 创建 Snowflake 实例
func NewSnowflakeWithNodeId(nodeId int64) *Snowflake {
	return NewSnowflake(nodeId, DEFAULT_CUSTOM_EPOCH)
}

// 自动生成 nodeId 创建 Snowflake 实例
func NewSnowflakeAuto() *Snowflake {
	nodeId := createNodeId()
	return NewSnowflake(nodeId, DEFAULT_CUSTOM_EPOCH)
}

// NextId 生成下一个全局唯一ID
// 在需要生成新的 ID 的时候，用当前时间戳加上当前时间戳内（也就是某一毫秒内）的计数器，拼接得到 UUID，
// 如果某一毫秒内的计数器被耗尽达到上限，会死循环直至这 1ms 过去
func (s *Snowflake) NextId() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTimestamp := s.timestamp()
	if currentTimestamp < s.lastTimestamp {
		// 系统时钟回退的情况
		panic("Invalid System Clock!")
	}

	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & int64(maxSequence)
		if s.sequence == 0 {
			// 序列号耗尽，等待下一毫秒
			currentTimestamp = s.waitNextMillis(currentTimestamp)
		}
	} else {
		// 新的毫秒，序列号从0开始
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	id := (currentTimestamp << (NODE_ID_BITS + SEQUENCE_BITS)) |
		(s.nodeId << SEQUENCE_BITS) |
		s.sequence

	return id
}

// 计算当前时间戳（毫秒）减去自定义纪元
func (s *Snowflake) timestamp() int64 {
	return time.Now().UnixMilli() - s.customEpoch
}

// 当序列耗尽时，等待下一毫秒
func (s *Snowflake) waitNextMillis(currentTimestamp int64) int64 {
	for currentTimestamp == s.lastTimestamp {
		currentTimestamp = s.timestamp()
	}
	return currentTimestamp
}

// 根据 MAC 地址生成 nodeId，如果失败则随机生成一个
func createNodeId() int64 {
	interfaces, err := net.Interfaces()
	if err == nil {
		var sb string
		for _, inter := range interfaces {
			mac := inter.HardwareAddr
			if len(mac) > 0 {
				for _, b := range mac {
					sb += fmt.Sprintf("%02X", b)
				}
			}
		}
		if len(sb) > 0 {
			hash := int64(hashString(sb))
			return hash & int64(maxNodeId)
		}
	}

	// 若无法获取MAC地址，则随机生成
	r, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return r.Int64() & int64(maxNodeId)
}

// 简单字符串哈希函数
func hashString(s string) int {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	return h
}

// 将生成的ID解析回 [timestamp, nodeId, sequence]
func (s *Snowflake) Parse(id int64) (timestamp int64, nodeId int64, sequence int64) {
	maskNodeId := ((int64(1) << NODE_ID_BITS) - 1) << SEQUENCE_BITS
	maskSequence := (int64(1) << SEQUENCE_BITS) - 1

	timestamp = (id >> (NODE_ID_BITS + SEQUENCE_BITS)) + s.customEpoch
	nodeId = (id & maskNodeId) >> SEQUENCE_BITS
	sequence = id & maskSequence

	return
}

func (s *Snowflake) String() string {
	return fmt.Sprintf("Snowflake Settings [EPOCH_BITS=%d, NODE_ID_BITS=%d, SEQUENCE_BITS=%d, CUSTOM_EPOCH=%d, NodeId=%d]",
		EPOCH_BITS, NODE_ID_BITS, SEQUENCE_BITS, s.customEpoch, s.nodeId)
}
