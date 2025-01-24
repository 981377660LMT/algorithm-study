// https://mp.weixin.qq.com/s?__biz=MzkxMjQzMjA0OQ==&mid=2247484405&idx=1&sn=ef9af767b617aca93bfaa9d01d6076c9
//
// 同步
// 异步
// 工程实践

package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	{
		observerA := NewBaseObserver("a")
		observerB := NewBaseObserver("b")
		observerC := NewBaseObserver("c")
		observerD := NewBaseObserver("d")

		sbus := NewBaseSubject()
		topic := "order_finish"
		sbus.Subscribe(topic, observerA)
		sbus.Subscribe(topic, observerB)
		sbus.Subscribe(topic, observerC)
		sbus.Subscribe(topic, observerD)

		sbus.Publish(context.Background(), &Event{
			Topic: topic,
			Val:   "order_id: xxx",
		})
	}

	{
		observerA := NewBaseObserver("a")
		observerB := NewBaseObserver("b")
		observerC := NewBaseObserver("c")
		observerD := NewBaseObserver("d")

		abus := NewAsyncSubject()
		defer abus.Stop()

		topic := "order_finish"
		abus.Subscribe(topic, observerA)
		abus.Subscribe(topic, observerB)
		abus.Subscribe(topic, observerC)
		abus.Subscribe(topic, observerD)

		abus.Publish(context.Background(), &Event{
			Topic: topic,
			Val:   "order_id: xxx",
		})

		<-time.After(time.Second)
	}
}

type Observable interface {
	Subscribe(topic string, observer Observer)
	Unsubscribe(topic string, observer Observer)
	Publish(ctx context.Context, e *Event)
}

type Observer interface {
	OnChange(ctx context.Context, e *Event) error
}

type Event struct {
	Topic string
	Val   any
}

func (e *Event) String() string {
	return e.Val.(string)
}

var _ Observer = (*BaseObserver)(nil)

type BaseObserver struct {
	name string
}

func NewBaseObserver(name string) *BaseObserver {
	return &BaseObserver{name: name}
}

func (o *BaseObserver) OnChange(ctx context.Context, e *Event) error {
	log.Printf("Observer %s received event %v", o.name, e)
	return nil
}

var _ Observable = (*BaseSubject)(nil)

type BaseSubject struct {
	mux       sync.RWMutex
	observers map[string]map[Observer]struct{}
}

func NewBaseSubject() *BaseSubject {
	return &BaseSubject{
		observers: make(map[string]map[Observer]struct{}),
	}
}

func (s *BaseSubject) Subscribe(topic string, observer Observer) {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.observers[topic]
	if !ok {
		s.observers[topic] = make(map[Observer]struct{})
	}
	s.observers[topic][observer] = struct{}{}
}

func (s *BaseSubject) Unsubscribe(topic string, observer Observer) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.observers[topic], observer)
}

// 同步模式
func (s *BaseSubject) Publish(ctx context.Context, e *Event) {
	s.mux.RLock()
	subscribers := s.observers[e.Topic]
	s.mux.RUnlock()

	errs := make(map[Observer]error)
	for observer := range subscribers {
		if err := observer.OnChange(ctx, e); err != nil {
			errs[observer] = err
		}
	}

	s.handleErrs(ctx, errs)
}

func (s *BaseSubject) handleErrs(_ context.Context, errs map[Observer]error) {
	for observer, err := range errs {
		//  处理 publish 失败的 observer
		log.Printf("Observer %s failed to handle event: %v", observer, err)
	}
}

// #region 异步模式

// 在异步模式下，会在 Subject 启动之初，异步启动一个守护协程，负责对接收到的错误进行后处理.
var _ Observable = (*AsyncSubject)(nil)

type observerWithErr struct {
	o   Observer
	err error
}

type AsyncSubject struct {
	BaseSubject
	errC chan *observerWithErr
	ctx  context.Context
	stop context.CancelFunc
}

func NewAsyncSubject() *AsyncSubject {
	res := &AsyncSubject{BaseSubject: *NewBaseSubject()}
	res.ctx, res.stop = context.WithCancel(context.Background())
	go res.handleErrs()
	return res
}

func (s *AsyncSubject) Stop() {
	s.stop()
}

func (s *AsyncSubject) Publish(ctx context.Context, e *Event) {
	s.mux.RLock()
	subscribers := s.observers[e.Topic]
	s.mux.RUnlock()

	for observer := range subscribers {
		// shadowing
		observer := observer
		go func() {
			if err := observer.OnChange(ctx, e); err != nil {
				// 正常运行时将错误发送到错误通道，当主上下文已关闭时不再发送错误
				select {
				case <-s.ctx.Done():
				case s.errC <- &observerWithErr{o: observer, err: err}:
				}
			}
		}()
	}
}

// 遇到的错误会通过 channel 统一汇总到 handleErr 的守护协程中进行处理.
func (s *AsyncSubject) handleErrs() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.errC:
			// 处理 publish 失败的 observer
			log.Printf("Observer %s failed to handle event: %v", err.o, err.err)
		}
	}
}

// #endregion
