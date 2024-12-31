// 单机版时间轮
// https://github.com/xiaoxuxiansheng/timewheel
// https://www.bilibili.com/video/BV1k8411r7E4/

package main

import (
	"container/list"
	"sync"
	"time"
)

type taskElement struct {
	task  func()
	pos   int    // 定时任务在环形数组所处的index
	cycle int    // 定时任务的延迟轮次. 指的是 curSlot 指针还要扫描过环状数组多少轮，才满足执行该任务的条件
	key   string // 全局唯一标识键
}

type TimeWheel struct {
	sync.Once

	interval time.Duration // 轮询时间间隔
	ticker   *time.Ticker  // 定时器

	stopCh       chan struct{}     // 关闭时间轮控制器
	addTaskCh    chan *taskElement // 添加任务的入口
	removeTaskCh chan string       // 删除任务的入口

	slots      []*list.List             // 时间轮环形数组(逻辑意义)
	curSlot    int                      // 当前时间轮指针
	keyToETask map[string]*list.Element // 任务 key 到任务结点的映射
}

// NewTimeWheel 创建一个时间轮
//
//	numSlot: 时间轮的槽数
//	interval: 时间轮的时间间隔
func NewTimeWheel(numSlot int, interval time.Duration) *TimeWheel {
	if numSlot <= 0 {
		numSlot = 10
	}
	if interval <= 0 {
		interval = time.Second
	}

	t := TimeWheel{
		interval:     interval,
		ticker:       time.NewTicker(interval),
		stopCh:       make(chan struct{}),
		addTaskCh:    make(chan *taskElement),
		removeTaskCh: make(chan string),
		slots:        make([]*list.List, 0, numSlot),
		keyToETask:   make(map[string]*list.Element),
	}
	for i := 0; i < numSlot; i++ {
		t.slots = append(t.slots, list.New())
	}

	// 异步启动时间轮常驻 goroutine
	go t.run()
	return &t
}

func (t *TimeWheel) Stop() {
	t.Do(func() {
		t.ticker.Stop()
		close(t.stopCh)
	})
}

func (t *TimeWheel) AddTask(key string, task func(), executeAt time.Time) {
	pos, cycle := t.getPosAndCircle(executeAt)
	t.addTaskCh <- &taskElement{
		pos:   pos,
		cycle: cycle,
		task:  task,
		key:   key,
	}
}

func (t *TimeWheel) RemoveTask(key string) {
	t.removeTaskCh <- key
}

func (t *TimeWheel) run() {
	defer func() {
		if err := recover(); err != nil {
			// ...
		}
	}()

	for {
		select {
		case <-t.stopCh:
			return
		case <-t.ticker.C:
			t.tick()
		case task := <-t.addTaskCh:
			t.addTask(task)
		case removeKey := <-t.removeTaskCh:
			t.removeTask(removeKey)
		}
	}
}

func (t *TimeWheel) tick() {
	list := t.slots[t.curSlot]
	defer t.circularIncr()
	t.execute(list)
}

func (t *TimeWheel) execute(l *list.List) {
	// 遍历每个 list
	for e := l.Front(); e != nil; {
		taskElement, _ := e.Value.(*taskElement)
		if taskElement.cycle > 0 {
			taskElement.cycle--
			e = e.Next()
			continue
		}

		// 执行任务
		go func() {
			defer func() {
				if err := recover(); err != nil {
					// ...
				}
			}()
			taskElement.task()
		}()

		// 执行任务后，从时间轮中删除
		next := e.Next()
		l.Remove(e)
		delete(t.keyToETask, taskElement.key)
		e = next
	}
}

func (t *TimeWheel) getPosAndCircle(executeAt time.Time) (int, int) {
	delay := int(time.Until(executeAt))
	cycle := delay / (len(t.slots) * int(t.interval))
	pos := (t.curSlot + delay/int(t.interval)) % len(t.slots)
	return pos, cycle
}

func (t *TimeWheel) addTask(task *taskElement) {
	list := t.slots[task.pos]
	if _, ok := t.keyToETask[task.key]; ok {
		t.removeTask(task.key)
	}
	eTask := list.PushBack(task)
	t.keyToETask[task.key] = eTask
}

func (t *TimeWheel) removeTask(key string) {
	eTask, ok := t.keyToETask[key]
	if !ok {
		return
	}
	delete(t.keyToETask, key)
	task, _ := eTask.Value.(*taskElement)
	_ = t.slots[task.pos].Remove(eTask)
}

func (t *TimeWheel) circularIncr() {
	t.curSlot = (t.curSlot + 1) % len(t.slots)
}
