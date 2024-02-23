package browser

import (
	"IM-Server/im/conn"
	_json "IM-Server/im/message/json"
	"encoding/json"
	"fmt"
	"sync"
)

func init() {
	recyclePool = sync.Pool{
		New: func() interface{} {
			return &readerRes{}
		},
	}

}

// recyclePool 回收池, 减少临时对象, 回收复用 readerRes
var recyclePool sync.Pool

type readerRes struct {
	err error
	m   *_json.ComMessage
}

// Recycle 回收当前对象, 一定要在用完后调用这个方法, 否则无法回收
func (r *readerRes) Recycle() {
	r.m = nil
	r.err = nil
	recyclePool.Put(r)
}

type defaultReader struct{}

func (d *defaultReader) ReadCh(conn conn.Connection) (<-chan *readerRes, chan<- interface{}) {
	c := make(chan *readerRes, 5)
	done := make(chan interface{})

	go func() {
		defer func() {
			e := recover()
			if e != nil {
				fmt.Println(e)
			}
		}()
		for {
			select {
			case <-done:
				goto CLOSE
			default:
				m, err := d.Read(conn)
				res := recyclePool.Get().(*readerRes)
				if err != nil {
					res.err = err
					c <- res
					goto CLOSE
				} else {
					res.m = m
					c <- res
				}
			}
		}
	CLOSE:
		close(c)
	}()
	return c, done
}

func (d *defaultReader) Read(conn conn.Connection) (*_json.ComMessage, error) {

	bytes, err := conn.Read()
	if err != nil {
		return nil, err
	}

	m := _json.NewEmptyMessage()
	err = json.Unmarshal(bytes, &m)
	return m, err
}
