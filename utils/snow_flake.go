package utils

import (
	"sync"
	"time"
)

const (
	epoch             = int64(1577808000000)                           // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits     = uint(37)                                       // 时间戳占用位数
	datacenteridBits  = uint(2)                                        // 数据中心id所占位数
	workeridBits      = uint(5)                                        // 机器id所占位数
	sequenceBits      = uint(5)                                        // 序列所占的位数
	timestampMax      = int64(-1 ^ (-1 << timestampBits))              // 时间戳最大值
	datacenteridMax   = int64(-1 ^ (-1 << datacenteridBits))           // 支持的最大数据中心id数量
	workeridMax       = int64(-1 ^ (-1 << workeridBits))               // 支持的最大机器id数量
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))               // 支持的最大序列id数量
	workeridShift     = sequenceBits                                   // 机器id左移位数
	datacenteridShift = sequenceBits + workeridBits                    // 数据中心id左移位数
	timestampShift    = sequenceBits + workeridBits + datacenteridBits // 时间戳左移位数
)

type Snowflake struct {
	sync.Mutex         // 锁
	Timestamp    int64 // 时间戳 ，毫秒
	Workerid     int64 // 工作节点
	Datacenterid int64 // 数据中心机房id
	Sequence     int64 // 序列号
}

func NewSnowflake() *Snowflake {
	return &Snowflake{
		Mutex:        sync.Mutex{},
		Timestamp:    time.Now().Unix(),
		Workerid:     1,
		Datacenterid: 1,
		Sequence:     5,
	}
}

func (s *Snowflake) NextVal() int64 {
	s.Lock()
	now := time.Now().UnixNano() / 1000000 // 转毫秒
	if s.Timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.Sequence = (s.Sequence + 1) & sequenceMask
		if s.Sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.Timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.Sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		return 0
	}
	s.Timestamp = now
	r := int64((t)<<timestampShift | (s.Datacenterid << datacenteridShift) | (s.Workerid << workeridShift) | (s.Sequence))
	s.Unlock()
	return r
}
