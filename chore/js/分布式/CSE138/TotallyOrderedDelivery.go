// totally_ordered_delivery.go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Message 消息结构
type Message struct {
	SenderID       int
	Content        string
	SequenceNumber int
}

// Sequencer 序列器结构
type Sequencer struct {
	incoming       chan Message
	receivers      []chan Message
	sequence       int
	receiversMutex sync.Mutex
	wg             *sync.WaitGroup
}

// NewSequencer 初始化序列器
func NewSequencer(wg *sync.WaitGroup) *Sequencer {
	return &Sequencer{
		incoming:  make(chan Message, 100),
		receivers: make([]chan Message, 0),
		sequence:  0,
		wg:        wg,
	}
}

// RegisterReceiver 注册一个接收者
func (s *Sequencer) RegisterReceiver(receiverChan chan Message) {
	s.receiversMutex.Lock()
	defer s.receiversMutex.Unlock()
	s.receivers = append(s.receivers, receiverChan)
}

// Start 启动序列器
func (s *Sequencer) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for msg := range s.incoming {
			s.sequence++
			msg.SequenceNumber = s.sequence
			fmt.Printf("Sequencer assigned Seq:%d to message from Sender:%d\n", s.sequence, msg.SenderID)
			s.broadcast(msg)
		}
		// 当incoming通道关闭时，关闭所有接收者通道
		s.closeAllReceivers()
	}()
}

// broadcast 将消息发送给所有注册的接收者
func (s *Sequencer) broadcast(msg Message) {
	s.receiversMutex.Lock()
	defer s.receiversMutex.Unlock()
	for _, receiverChan := range s.receivers {
		// 非阻塞发送，避免一个接收者阻塞所有发送
		select {
		case receiverChan <- msg:
		default:
			fmt.Printf("Warning: Receiver channel is full. Dropping message Seq:%d\n", msg.SequenceNumber)
		}
	}
}

// closeAllReceivers 关闭所有接收者通道
func (s *Sequencer) closeAllReceivers() {
	s.receiversMutex.Lock()
	defer s.receiversMutex.Unlock()
	for _, receiverChan := range s.receivers {
		close(receiverChan)
	}
}

// Sender 发送者结构
type Sender struct {
	SenderID          int
	sequencerIncoming chan Message
	messages          []string
	wg                *sync.WaitGroup
}

// NewSender 初始化发送者
func NewSender(senderID int, sequencerIncoming chan Message, messages []string, wg *sync.WaitGroup) *Sender {
	return &Sender{
		SenderID:          senderID,
		sequencerIncoming: sequencerIncoming,
		messages:          messages,
		wg:                wg,
	}
}

// Start 启动发送者
func (s *Sender) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for _, msg := range s.messages {
			message := Message{
				SenderID: s.SenderID,
				Content:  msg,
			}
			fmt.Printf("Sender %d sending message: '%s'\n", s.SenderID, msg)
			s.sequencerIncoming <- message
			time.Sleep(100 * time.Millisecond) // 模拟发送间隔
		}
	}()
}

// Receiver 接收者结构
type Receiver struct {
	ReceiverID        int
	sequencerOutgoing chan Message
	receivedMessages  []Message
	wg                *sync.WaitGroup
}

// NewReceiver 初始化接收者
func NewReceiver(receiverID int, sequencerOutgoing chan Message, wg *sync.WaitGroup) *Receiver {
	return &Receiver{
		ReceiverID:        receiverID,
		sequencerOutgoing: sequencerOutgoing,
		receivedMessages:  make([]Message, 0),
		wg:                wg,
	}
}

// Start 启动接收者
func (r *Receiver) Start() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for msg := range r.sequencerOutgoing {
			r.receivedMessages = append(r.receivedMessages, msg)
			fmt.Printf("Receiver %d received: %s\n", r.ReceiverID, msg.Content)
			time.Sleep(50 * time.Millisecond) // 模拟处理时间
		}
	}()
}

func main() {
	var wg sync.WaitGroup

	// 创建序列器
	sequencer := NewSequencer(&wg)
	sequencer.Start()

	// 创建接收者
	receiverCount := 2
	receivers := make([]*Receiver, 0, receiverCount)
	for i := 1; i <= receiverCount; i++ {
		receiverChan := make(chan Message, 100) // 为每个接收者创建独立的通道
		sequencer.RegisterReceiver(receiverChan)
		receiver := NewReceiver(i, receiverChan, &wg)
		receiver.Start()
		receivers = append(receivers, receiver)
	}

	// 创建发送者
	sender1 := NewSender(1, sequencer.incoming, []string{"A1", "A2", "A3"}, &wg)
	sender2 := NewSender(2, sequencer.incoming, []string{"B1", "B2", "B3"}, &wg)

	sender1.Start()
	sender2.Start()

	// 等待发送者完成发送
	wg.Wait()

	// 关闭序列器的incoming通道，触发序列器关闭接收者通道
	close(sequencer.incoming)

	// 等待所有goroutine完成
	wg.Wait()

	// 打印接收者收到的消息
	fmt.Println("\nReceiver1 Messages:")
	for _, msg := range receivers[0].receivedMessages {
		fmt.Printf("Seq:%d | Sender:%d | Content:%s\n", msg.SequenceNumber, msg.SenderID, msg.Content)
	}

	fmt.Println("\nReceiver2 Messages:")
	for _, msg := range receivers[1].receivedMessages {
		fmt.Printf("Seq:%d | Sender:%d | Content:%s\n", msg.SequenceNumber, msg.SenderID, msg.Content)
	}
}
