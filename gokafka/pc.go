package gokafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"strconv"
	"time"
)

type PC struct {
	blockingQueue chan interface{}
	batchSize int
	batchTimeOut time.Duration
}

func (pc *PC) Init(capacity int, batchSize int, batchTimeOut time.Duration){
	fmt.Printf("capacity %d, batchSize %d, batchTimeOut %d\n",capacity,batchSize,batchTimeOut)
	pc.blockingQueue = make(chan interface{}, capacity)
	pc.batchSize=batchSize
	pc.batchTimeOut=batchTimeOut
}

func (pc *PC) Produce(object interface{})  {
	pc.blockingQueue <- object
}

func (pc *PC) Subscribe(mapTo func(interface{})kafka.Message, consume func([]kafka.Message))  {
	for {
		L := make([]kafka.Message, pc.batchSize)
		size := 0
		run := true
		timer := time.NewTimer(time.Hour * 999999)//inf
		startTimer:=false
		for run {
			select {
			case x := <- pc.blockingQueue:
				L[size] = mapTo(x)
				size++
				if size >= pc.batchSize {//batchSize
					//fmt.Println("buffer is full, send to kafka")
					run = false
				} else if !startTimer {
					startTimer = true
					timer.Reset(pc.batchTimeOut)//batchTimeout
				}
			case <- timer.C:
				//fmt.Println("it's time to send to kafka, batch size "+strconv.Itoa(size))
				run=false
			}
		}
		if size != len(L) {
			L = L[:size]
		}
		consume(L)
	}
}
