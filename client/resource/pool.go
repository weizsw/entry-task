package resource

import (
	"errors"
	"log"
	"net"
	"sync"
)

var CP *Pool

type Pool struct {
	m        sync.Mutex
	resource chan net.Conn
	maxSize  int
	usedSize int
	factory  func() (net.Conn, error)
	closed   bool
}

func NewConnPool(factory func() (net.Conn, error), cap int) (*Pool, error) {
	if cap <= 0 {
		return nil, errors.New("cap could not be zero")
	}

	cp := &Pool{
		m:        sync.Mutex{},
		resource: make(chan net.Conn, cap),
		maxSize:  cap,
		usedSize: 0,
		factory:  factory,
		closed:   false,
	}

	for i := 0; i < cap; i++ {
		connRes, err := cp.factory()
		if err != nil {
			return nil, errors.New("factory err")
		}

		cp.resource <- connRes
	}

	return cp, nil
}

func (p *Pool) Get() (net.Conn, error) {
	p.m.Lock()
	defer p.m.Unlock()

	select {
	case r, ok := <-p.resource:
		if !ok {
			log.Println("Pool has been closed.")
			return nil, errors.New("Pool has been closed.")
		}
		p.usedSize++
		return r, nil
	default:
		if p.usedSize < p.maxSize {
			log.Printf("Acquire:"+"New Resource."+
				"resource present size/max: %d/%d\n", p.usedSize, p.maxSize)
			p.usedSize++
			return p.factory()
		} else {
			log.Println("create err")
			return nil, nil
		}
	}
}

func (p *Pool) Put(r net.Conn) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		r.Close()
		return
	}

	p.usedSize--

	select {
	case p.resource <- r:
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}
