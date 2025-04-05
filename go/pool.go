package main

import "sync"

type Pool[T struct{}] struct {
	items []*T
	mutex sync.Mutex
}

func (p *Pool[T]) Get() *T {
	p.mutex.Lock()

	poolLen := len(p.items)
	var data *T

	if poolLen == 0 {
		data = &T{}
		p.items = append(p.items, data)
	} else {
		data = p.items[poolLen-1]
		p.items = p.items[:poolLen-1]
	}
	p.mutex.Unlock()
	return data
}

func (p *Pool[T]) Release(d *T) {
	p.mutex.Lock()
	p.items = append(p.items, d)
	p.mutex.Unlock()
}
