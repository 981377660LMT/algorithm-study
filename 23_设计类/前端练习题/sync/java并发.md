# Java 并发工具详解

Java 提供了丰富的并发工具，主要通过 `java.util.concurrent` 包及其子包提供。以下是对 Java 并发工具的全面介绍：

## 一、线程基础工具

### 1. Thread 和 Runnable

- **Thread**: Java 中线程的基本实现类
- **Runnable**: 定义任务的接口，避免继承 Thread 的限制
- **Callable**: 带返回值的任务接口

### 2. 线程生命周期管理

- **Thread.sleep()**: 使线程休眠指定时间
- **Thread.yield()**: 提示调度器可以切换到其他线程
- **Thread.join()**: 等待指定线程结束
- **Thread.interrupt()**: 中断线程

## 二、同步工具类

### 1. 基本同步工具

- **synchronized**: 关键字，基于对象监视器实现的同步机制
- **volatile**: 关键字，保证变量的可见性和有序性
- **wait()/notify()/notifyAll()**: Object 类提供的线程通信机制
- **ReentrantLock**: 可重入锁，synchronized 的替代方案，提供更多功能

### 2. `java.util.concurrent.locks` 包

- **Lock 接口**: 比 synchronized 更灵活的锁机制

  - **ReentrantLock**: 可重入互斥锁，支持公平和非公平锁定
  - **ReadWriteLock 接口**: 读写锁接口
    - **ReentrantReadWriteLock**: 可重入读写锁实现
  - **StampedLock**: Java 8 引入，支持乐观读的读写锁

- **Condition 接口**: 替代 Object 的 wait/notify 机制
  - 支持多个等待集(wait-set)
  - 提供 await/signal/signalAll 方法

### 3. `java.util.concurrent` 同步工具类

- **CountDownLatch**: 允许线程等待一个或多个线程完成操作
- **CyclicBarrier**: 让一组线程在某个点上同步，可重用
- **Phaser**: Java 7 引入，比 CyclicBarrier 更灵活的同步屏障
- **Semaphore**: 控制并发访问的数量
- **Exchanger**: 两个线程在某点交换对象

## 三、并发容器

### 1. 集合类

- **ConcurrentHashMap**: 线程安全的 HashMap
- **ConcurrentSkipListMap**: 线程安全的有序 map
- **ConcurrentSkipListSet**: 线程安全的有序 set
- **CopyOnWriteArrayList**: 线程安全的 ArrayList，适合读多写少场景
- **CopyOnWriteArraySet**: 基于 CopyOnWriteArrayList 的线程安全 Set

### 2. 队列

- **BlockingQueue 接口**: 提供阻塞的入队和出队操作

  - **ArrayBlockingQueue**: 基于数组的有界阻塞队列
  - **LinkedBlockingQueue**: 基于链表的可选有界阻塞队列
  - **PriorityBlockingQueue**: 基于优先级的无界阻塞队列
  - **DelayQueue**: 延迟元素出队的无界阻塞队列
  - **SynchronousQueue**: 没有内部容量的阻塞队列
  - **LinkedTransferQueue**: Java 7 引入，结合 SynchronousQueue 和 LinkedBlockingQueue 功能

- **BlockingDeque 接口**: 双端阻塞队列接口
  - **LinkedBlockingDeque**: 基于链表的可选有界阻塞双端队列

### 3. 其他并发容器

- **ConcurrentLinkedQueue**: 非阻塞的线程安全队列
- **ConcurrentLinkedDeque**: 非阻塞的线程安全双端队列

## 四、执行器框架 (Executor Framework)

### 1. 基础接口

- **Executor**: 执行提交的 Runnable 任务的对象
- **ExecutorService**: 扩展 Executor，管理任务终止的完整服务
- **ScheduledExecutorService**: 支持定时和周期任务的 ExecutorService

### 2. 线程池实现

- **Executors 工厂类**:

  - **newFixedThreadPool()**: 固定大小线程池
  - **newCachedThreadPool()**: 可缓存的线程池
  - **newSingleThreadExecutor()**: 单线程执行器
  - **newScheduledThreadPool()**: 支持定时任务的线程池
  - **newWorkStealingPool()**: Java 8 引入，工作窃取线程池

- **ThreadPoolExecutor**: 高度可配置的线程池实现
- **ScheduledThreadPoolExecutor**: 支持定时任务的线程池实现
- **ForkJoinPool**: Java 7 引入，支持工作窃取算法的线程池

### 3. 任务提交和结果

- **Future 接口**: 表示异步计算的结果
- **CompletableFuture**: Java 8 引入，支持函数式编程的 Future
- **FutureTask**: Future 的标准实现类

## 五、原子变量类 (Atomic Variables)

### 1. 基本类型原子类

- **AtomicBoolean**: 原子布尔值
- **AtomicInteger**: 原子整数
- **AtomicLong**: 原子长整数

### 2. 引用类型原子类

- **AtomicReference<V>**: 原子引用
- **AtomicMarkableReference<V>**: 带标记位的原子引用
- **AtomicStampedReference<V>**: 带版本号的原子引用

### 3. 数组类型原子类

- **AtomicIntegerArray**: 原子整数数组
- **AtomicLongArray**: 原子长整数数组
- **AtomicReferenceArray<V>**: 原子引用数组

### 4. 字段更新器

- **AtomicIntegerFieldUpdater<T>**: 整数字段的原子更新器
- **AtomicLongFieldUpdater<T>**: 长整数字段的原子更新器
- **AtomicReferenceFieldUpdater<T,V>**: 引用字段的原子更新器

### 5. Java 8 新增原子类

- **DoubleAdder**: 高并发下性能更好的累加器
- **LongAdder**: 高并发下性能更好的累加器
- **DoubleAccumulator**: 支持自定义操作的原子累加器
- **LongAccumulator**: 支持自定义操作的原子累加器

## 六、并发工具类

### 1. 时间工具

- **TimeUnit 枚举**: 提供时间单位转换
  - NANOSECONDS, MICROSECONDS, MILLISECONDS, SECONDS, MINUTES, HOURS, DAYS

### 2. 同步工具

- **CountedCompleter**: Fork/Join 框架中的任务，完成后自动触发回调
- **CompletionService 接口**: 管理异步任务执行和结果获取
  - **ExecutorCompletionService**: 标准实现类

### 3. Java 9+ 新工具

- **Flow API**: Java 9 引入的响应式流接口
- **SubmissionPublisher**: Flow.Publisher 的标准实现类

## 七、线程安全集合适配器

### 1. 集合同步包装器

- **Collections.synchronizedCollection()**
- **Collections.synchronizedList()**
- **Collections.synchronizedMap()**
- **Collections.synchronizedSet()**
- **Collections.synchronizedSortedMap()**
- **Collections.synchronizedSortedSet()**

## 八、并发设计模式支持

### 1. 不变性支持

- **final 关键字**: 创建不可变对象
- **ImmutableList/Set/Map** (Guava 库): 不可变集合

### 2. 线程封闭

- **ThreadLocal<T>**: 提供线程局部变量
- **InheritableThreadLocal<T>**: 可继承的线程局部变量

## 九、Fork/Join 框架

- **ForkJoinPool**: 支持 Fork/Join 任务的线程池
- **ForkJoinTask<V>**: Fork/Join 任务抽象类
  - **RecursiveTask<V>**: 有返回值的递归任务
  - **RecursiveAction**: 无返回值的递归任务
  - **CountedCompleter<T>**: 完成执行后触发回调的任务

## 十、Java 内存模型

- **Happens-Before 规则**: 保证并发操作的可见性和有序性
- **内存屏障**: 防止指令重排序
- **JMM (Java Memory Model)**: 定义线程安全的规范

## 十一、并发编程最佳实践工具

### 1. 锁优化

- **偏向锁(Biased Locking)**: JVM 优化技术
- **轻量级锁(Lightweight Locking)**: JVM 优化技术
- **自旋锁(Spin Locking)**: JVM 优化技术
- **锁粗化(Lock Coarsening)**: JVM 优化技术
- **锁消除(Lock Elision)**: JVM 优化技术

### 2. 线程安全性标注

- **@ThreadSafe**: 标注线程安全的类(来自JSR-305)
- **@NotThreadSafe**: 标注非线程安全的类(来自JSR-305)
- **@Immutable**: 标注不可变类(来自JSR-305)
- **@GuardedBy**: 标注受特定锁保护的字段(来自JSR-305)

### 3. 调试和监控工具

- **jconsole**: JVM 监控和管理控制台
- **jvisualvm**: Java 应用程序可视化监控工具
- **jstack**: 生成线程堆栈跟踪
- **Thread.getAllStackTraces()**: 编程方式获取所有线程栈跟踪

## 十二、Java 9+ 的新增特性

### 1. Java 9

- **CompletableFuture 增强**: 添加超时、延迟和其他组合操作
- **Flow API**: 响应式流支持

### 2. Java 10+

- **var 支持**: 局部变量类型推断，简化并发代码

### 3. Java 19+

- **结构化并发(Preview)**: 简化复杂的异步编程
- **虚拟线程(Preview)**: 轻量级线程实现

## 十三、第三方并发库

### 1. Guava

- **ListenableFuture**: 可监听的 Future
- **MoreExecutors**: 更多 Executor 工具
- **RateLimiter**: 速率限制器

### 2. Reactive Streams 实现

- **RxJava**: 响应式编程库
- **Project Reactor**: Spring 生态系统的响应式库
- **Akka Streams**: 基于 Actor 模型的流处理库

### 3. 其他

- **Disruptor**: 高性能线程间消息传递框架
- **JCTools**: Java 并发工具，提供高性能非阻塞数据结构

## 总结

Java 提供了全面的并发工具集，从低级别的同步原语到高级别的执行器框架，使开发者能够构建高效、可靠的并发应用程序。这些工具通过抽象复杂的底层机制，显著降低了编写并发代码的难度和风险。

随着 Java 版本的迭代，并发工具库也在不断增强和完善，引入了更多现代编程模式的支持，如响应式编程、函数式编程和结构化并发等，使 Java 在并发编程领域保持竞争力。
