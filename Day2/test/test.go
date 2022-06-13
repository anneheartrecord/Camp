package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	m    = make(map[int]int, 0)
	lock sync.Mutex
)

func addWithLock() {
	for i := 0; i < 2000; i++ {
		lock.Lock()
		m[i] += 1
		lock.Unlock()
	}
}
func addWithoutLock() {
	for i := 0; i < 2000; i++ {
		m[i] += 1
	}
}

func main() {
	// slice map is not thread safe
	// but variate and chan is thread safe
	// so  we  can use different goroutine to change a variate
	for i := 0; i < 5; i++ {
		go addWithLock()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("addwithlock", m)

	for i := 0; i < 5; i++ {
		go addWithoutLock()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("addwithoutlock", m)
}
