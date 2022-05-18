package resource

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var ConnPool *Pool

type Pool struct {
	m           sync.Mutex
	conns       chan net.Conn
	factory     func() (net.Conn, error)
	closed      bool
	connTimeOut time.Duration
}

type Conn struct {
	conn net.Conn
	time time.Time
}

func NewConnPool(factory func() (net.Conn, error), cap int, connTimeOut time.Duration) (*Pool, error) {
	if cap <= 0 {
		return nil, errors.New("cap could not be zero")
	}

	if connTimeOut < 0 {
		return nil, errors.New("connTimeOut could not be negative")
	}

	cp := &Pool{
		m:           sync.Mutex{},
		conns:       make(chan net.Conn, cap),
		factory:     factory,
		closed:      false,
		connTimeOut: connTimeOut,
	}

	for i := 0; i < cap; i++ {
		res, err := cp.factory()
		if err != nil {
			cp.Close()
			return nil, errors.New("factory err")
		}

		cp.conns <- res
	}

	return cp, nil
}

func (p *Pool) Get() (net.Conn, error) {
	if p.closed {
		return nil, errors.New("pool closed")
	}
	timeout := time.After(p.connTimeOut)
	for {
		select {
		case res, ok := <-p.conns:
			{
				if !ok {
					return nil, errors.New("pool closed")
				}

				return res, nil
			}
		case <-timeout:
			return nil, errors.New("timeout")
		default:
			{
				log.Println("making new")
				res, err := p.factory()
				if err != nil {
					log.Println("new failed")
					return nil, err
				}
				return res, nil
			}
		}
	}
}

func (p *Pool) Put(r net.Conn) {
	if p.closed {
		return
	}

	select {
	case p.conns <- r:

	default:
		log.Println("put:", "closing")
		r.Close()
	}
}

func (cp *Pool) Close() {
	if cp.closed {
		return
	}
	cp.m.Lock()
	cp.closed = true

	close(cp.conns)
	for conn := range cp.conns {
		conn.Close()
	}
	cp.m.Unlock()
}

func (p *Pool) Len() int {
	return len(p.conns)
}
