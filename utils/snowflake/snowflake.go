package snowflake

import (
	"errors"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime   int64 = 1525705533000 // 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

var work *Worker
var works = make(map[int64]*Worker)

func InitSnowFlake(workerId int64) error {
	if workerId < 0 || workerId > workerMax {
		return errors.New("Worker ID excess of quantity")
	}
	// 生成一个新节点
	work = &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}
	return nil
}
func GetNextIdByServer(workerId int64) int64 {
	work, h := works[workerId]
	if !h {
		work = &Worker{
			timestamp: 0,
			workerId:  workerId,
			number:    0,
		}
		works[workerId] = work
	}
	return getNextWorkId(work)
}
func GetNextID() int64 {
	if work == nil {
		glog.Fatalf(gctx.New(), "雪花算法未初始化")
	}
	return getNextWorkId(work)
}
func getNextWorkId(work *Worker) int64 {
	work.mu.Lock()
	defer work.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if work.timestamp == now {
		work.number++
		if work.number > numberMax {
			for now <= work.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		work.number = 0
		work.timestamp = now
	}
	ID := int64((now-startTime)<<timeShift | (work.workerId << workerShift) | (work.number))
	return ID
}
