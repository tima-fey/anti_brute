package localDB

import (
	"time"
)

type Bucket interface {
	Add(string, out chan<- bool)
	Clear(string)
	ClearAll()
}

type BaseBucket struct {
	backetSize int
	interval   int64
	timestamps map[string][]int64
}

type BaseDB struct {
	Address  BaseBucket
	Login    BaseBucket
	Password BaseBucket
}

func DbInit() BaseDB {
	DB := BaseDB{BaseBucket{}, BaseBucket{}, BaseBucket{}}
	DB.Address.backetSize = 10
	DB.Address.interval = 60
	DB.Address.timestamps = make(map[string][]int64)
	DB.Login.backetSize = 10
	DB.Login.interval = 60
	DB.Login.timestamps = make(map[string][]int64)
	DB.Password.backetSize = 10
	DB.Password.interval = 60
	DB.Password.timestamps = make(map[string][]int64)
	return DB
}

func (b *BaseBucket) Add(k string, out chan<- bool) {
	currentTime := time.Now().Unix()
	if len(b.timestamps[k]) >= b.backetSize {
		for i, time := range b.timestamps[k] {
			if time < currentTime-b.interval {
				if i < len(b.timestamps[k])-1 {
					b.timestamps[k] = b.timestamps[k][i+1:]
				} else {
					b.timestamps[k] = nil
				}
			} else {
				break
			}
		}
		if len(b.timestamps) >= b.backetSize {
			out <- false
			return
		}
	}
	b.timestamps[k] = append(b.timestamps[k], currentTime)
	out <- true
}

func (b *BaseBucket) Clear(k string) {
	b.timestamps[k] = nil
}
func (b *BaseBucket) ClearAll() {
	b.timestamps = nil
}
