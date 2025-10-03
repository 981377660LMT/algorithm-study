# Golang pprof 案例实战与原理解析

https://mp.weixin.qq.com/s/Qwmo9FHCF010-0rMUbyuww

- profile：探测各函数对 cpu 的占用情况

cpu 分析是在一段时间内进行打点采样，通过查看采样点在各个函数栈中的分布比例，以此来反映各函数对 cpu 的占用情况.

- heap：探测内存分配情况
- block：探测阻塞情况 （包括 mutex、chan、cond 等）

查看某个 goroutine 陷入 waiting 状态（被动阻塞，通常因 gopark 操作触发，比如因加锁、读 chan 条件不满足而陷入阻塞）的触发次数和持续时长.

- mutex：探测互斥锁占用情况
- goroutine：探测协程使用情况

---

- Cpu 分析：   
  启动定时器 timer，定期向各个 thread 发送 SIGPROF 信号    
  处理 SIGPROF 信号时记录函数栈信息，通过这种抽样的方式反映各个函数对 CPU 占用情况
- Heap 分析：   
  每分配指定大小的内存，就会采样一笔内存分配信息并记录在全局变量 memBlock 中    
  每轮 gc 结束前，同样在 memBlock 中记录内存释放信息    
  读取内存指标时，遍历各个 memBucket 加载结果返回
- Block 分析：   
  根据传参确定采集频率，在 goroutine 阻塞并被重新唤醒后进行阻塞信息上报，将其存储在全局变量 blockBucket 中    
  读取 block 指标时，遍历各个 blockBucket 加载结果返回
- Mutex 分析：   
  根据传参确定采集频率，在解锁前上报加锁时长信息，存储在全局变量 mutexBucket 中    
  读取 mutex 指标时，遍历各个 mutexBucket 加载结果返回
- Goroutine 分析：   
  读取 goroutine 指标时遍历各个 g，获取其栈信息后返回
