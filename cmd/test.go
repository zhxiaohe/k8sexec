package main

import (
	"fmt"
	"time"
)

type server struct {
	r chan string
	w chan string
	e chan string
}

func main1() {
	test := init1()
	fmt.Println(test)
	time.Sleep(5 * time.Second)
}

func init1() string {
	r := make(chan string)
	w := make(chan string)
	e := make(chan string)
	s := &server{r, w, e}
	go s.read()
	go s.write()
	go client(s)

	return "run"
}
func client(s *server) {
	for i := 0; i < 60; i++ {
		str := fmt.Sprintf("hi %d", i)
		s.r <- str
		time.Sleep(1 * time.Second)
	}
}

func (s server) read() {
	for {
		select {
		case recv := <-s.r:
			fmt.Println("read: ", recv)
			s.w <- recv + "+ack"
		}
	}
}

func (s server) write() {
	for {
		select {
		case ack := <-s.w:
			fmt.Println("write: ", ack)
		}
	}
}
