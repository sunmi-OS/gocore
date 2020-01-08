package gokafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)
// 用于写kafka阻塞模式，这样当写kafka失败时，可以收到该错误，并进行相应的日志记录或者处理，
// 一次传输较多消息给kafka这样具有较高的吞吐，所以采用生产者消费者模式

type PC struct {
	blockingQueue chan interface{}
	batchSize int						// 一个批次的最大消息数量
	batchTimeOut time.Duration			// 准备一个批次的消息花费的最长时间
	shutdown chan int
}

// 设置阻塞队列的容量，批次大小，批次超时时间
// capacity     阻塞队列的容量，建议这个值设置为 batchSize 的 4 倍
// batchSize    批次大小
// batchTimeOut 批次超时时间
func (pc *PC) Init(capacity int, batchSize int, batchTimeOut time.Duration){
	fmt.Printf("capacity %d, batchSize %d, batchTimeOut %d\n",capacity,batchSize,batchTimeOut)
	pc.blockingQueue = make(chan interface{}, capacity)
	pc.batchSize=batchSize
	pc.batchTimeOut=batchTimeOut
	pc.shutdown = make(chan int)
}

// 向阻塞队列中生产消息，当阻塞队列已经满时，会阻塞
func (pc *PC) Produce(object interface{})  {
	pc.blockingQueue <- object
}

// 消费阻塞队列中的数据
// mapTo   由于向队列中写数据是一个 interface{} 对象，你需要在这个回调中实现你自己的序列方式，将其传换成 kafka.Message 对象
// consume 每当有可消费的批次时，该方法就会被回调，你可以在这里实现你自己的消费逻辑，通常这里使用gokafka.Producer.ProduceMsgs 并将消息阻塞的写入kafka
// 下面是一个使用例子
//	pc.Subscribe(func(x interface{}) kafka.Message {
//		return x.(kafka.Message)
//	}, func(messages []kafka.Message) {
//		err := producer.ProduceMsgs(messages)
//		if err != nil {
//			log.Error(err)
//		} else {
//			fmt.Printf("Success write `%d` to kafka\n", len(messages))
//		}
//	})
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
				run = false
			case <- pc.shutdown:
				// add shutdown hook
				tmp := make([]kafka.Message, len(pc.blockingQueue))
				for i := len(tmp) - 1; i >= 0; i-- {
					tmp[i] = mapTo(<- pc.blockingQueue)
				}
				messages := append(L[:size], tmp...)
				fmt.Printf("try flush buffer size %d\n", len(messages))
				consume(messages)
				pc.shutdown <- 1 //send successful signal
				return
			}
		}
		if size != len(L) {
			L = L[:size]
		}
		consume(L)
	}
}

//cancel the produce, blocking until the kafka cache buffer is successfully refreshed
func (pc *PC)Cancel() bool {
	pc.shutdown <- 0 // send shutdown signal
	_, ok := <- pc.shutdown // wait for shutdown
	return ok
}
