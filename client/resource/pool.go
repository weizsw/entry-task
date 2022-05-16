package resource

import (
	"errors"
	"io"
	"log"
	"sync"
)

var CP *Pool

type Pool struct {
	m        sync.Mutex
	resource chan io.Closer
	maxSize  int
	usedSize int
	factory  func() (io.Closer, error)
	closed   bool
}

func NewConnPool(factory func() (io.Closer, error), cap int) (*Pool, error) {
	if cap <= 0 {
		return nil, errors.New("cap could not be zero")
	}

	cp := &Pool{
		m:        sync.Mutex{},
		resource: make(chan io.Closer, cap),
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

func (p *Pool) Get() (io.Closer, error) {
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

func (p *Pool) Put(r io.Closer) {
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
