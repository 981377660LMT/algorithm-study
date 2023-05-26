1. 跨平台性是 JVM 存在的最大的亮点:如何实现平台无关
   不同平台的 JVM 解析字节码 转成具体的机器指令
   查看字节码：javac -p 反编译
2. 我们利用 JDK 开发了属于我们自己的程序，通过 JDK 的 javac 工具包进行了编译，将 Java 文件编译成为了 class 文件（字节码文件），在 JRE 上运行这些文件的时候，JVM 进行了这些文件（字节码文件）的翻译，翻译给操作系统，映射到 CPU 指令集或者是操作系统调用，最终完成了我们的代码程序的顺利运行。
3. JVM 如何加载并执行.class 文件
   class loader 将字节码准尉 JVM 中的 Class<?> 对象
   runtime data areas 主要包括五个小模块，Java 堆， Java 栈，本地方法栈，方法区，寄存器
   excution engine
   native interface
4. classLoader 的双亲委派机制 避免多份同样字节码的加载

verbose 口语化

4. 内存模型
   堆栈区别
   - 栈自动释放 堆需要 GC
   - 栈比堆小
   - 栈产生的碎片远小于堆
   - 栈支持静态和动态分配，而堆只支持动态分配
   - 栈效率比堆高
5. 垃圾**标记**算法
   被判断为垃圾:没有引用
   判断方法:引用计数(缺点是循环引用)/可达性算法(GC Root 开始遍历)
6. 垃圾**回收**算法
   **复制算法**(Copying)

   - 存活的对象从**对象面**被复制到**空闲面** 然后清除对象面对象内存
     适用于存活率低的情况(回收**年轻代**)
     解决碎片化问题，顺序分配内存，简单高效
     使用 From -> 触发 GC 标记整理 -> 拷贝到 To -> 回收 From -> 名称互换重复之前

     **标记清除**(Mark and Sweep)

   - 可达性算法标记存活对象=>回收不可达对象内存
     导致碎片多

   **标记整理**(Compacting)

   - 可达性算法标记存活对象=>移动存活对象，按照内存地址排列，然后回收末端内存地址之后的全部内存
     标记清除进化版,解决碎片化问题,适用于存活率高的情景，适用于**老年代**

   **分代收集算法**(Generational Collector)

   - 垃圾回收算法组合拳,按对象生命周期采用不同策略 (新空间：老空间=1:2)
     **年轻代**(eden 区生成，Survivor 区的 from 和 to 区交换)：复制算法 (Minor GC),空间换时间
     **老年代**：标记清除、标记整理 (老年代空间不足时 Full GC)；标记是三色并发标记;老生代区域垃圾回收不适合复制算法，老生代**空间大**,一分为二，会造成一半的空间浪费，存放数据多复制时间长。

7. 晋升:晋升就是将新生代对象移动至老生代。
   什么时候触发晋升操作?

   一轮 GC 之后还存活的新生代对象需要晋升
   在拷贝过程中，To 空间的使用率超过 25%，将这次的活动对象都移动至老生代空间

8. v8/jvm/go 三色标记过程
   三色标记使得垃圾回收可以分片并发(concurrent),可以暂停重启
   ![三色](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/23aea3d2c7ea43b6b9f39c45ddf68499~tplv-k3u1fbpfcp-watermark.awebp)
   白色：表示对象尚未被垃圾收集器访问过(不可达。)
   黑色：表示对象已经被垃圾收集器访问过，且这个对象的所有引用都已经扫描过。(可达)
   灰色：表示对象已经被垃圾收集器访问过，但这个对象上至少存在一个引用还没有被扫描过。(待处理)

   1. 在 GC 并发开始的时候，所有的对象均为白色；
   2. 在将所有的 GC Roots 直接应用的对象标记为灰色集合；
   3. 如果判断灰色集合中的对象不存在子引用，则将其放入黑色集合，若存在子引用对象，则将其所有的子引用对象存放到灰色集合，当前对象放入黑色集合
   4. 按照此步骤 3 ，依此类推，直至灰色集合中所有的对象变黑后，本轮标记完成，并且在白色集合内的对象称为不可达对象，即垃圾对象。
   5. 标记结束后，为白色的对象为 GC Roots 不可达，可以进行垃圾回收。

      **综合来说就是树遍历中实际存在三种状态：已经遍历完的，遍历中的，没遍历的**

      而当需要支持并发标记时，即标记期间应用线程还在继续跑，对象间的引用可能发生变化，多标和漏标的情况就有可能发生。

      多标-浮动垃圾：黑 D 指向灰 E 的引用断开,**因为 E 已经变为灰色了，其仍会被当作存活对象继续遍历下去。**最终的结果是：这部分对象仍会被标记为存活，即本轮 GC 不会回收这部分内存。这部分本应该回收 但是没有回收到的内存，被称之为“浮动垃圾”。浮动垃圾并不会影响应用程序的正确性，只是需要等到下一轮垃圾回收中才被清除。
      ![img](https://img2020.cnblogs.com/blog/1153954/202012/1153954-20201220213532679-409632912.png)

      漏标-读写屏障:不能让黑色节点指向白色节点！每当发生引用变化时，需要立刻对被引用节点进行着色：即白的立刻染灰，灰的和黑的不变。
      ![漏标](https://img2020.cnblogs.com/blog/1153954/202012/1153954-20201220214557708-22631031.png)

      ```JS
      var G = objE.fieldG; // 1.读
      objE.fieldG = null;  // 2.写
      objD.fieldG = G;     // 3.写
      写屏障用于拦截第二和第三步；而读屏障则是拦截第一步。
      它们的拦截的目的很简单：就是在读写前后，将对象 G 给记录下来。
      ```

9. 强引用/软引用/弱引用/虚引用
   强引用

```JAVA
Object obj = new Object()
```

弱引用 GC 会被回收

```JAVA
// JAVA
```

11. 串行=>批处理=>进程=>线程
    进程是资源分配的最小单位 线程是 CPU 调度的最小单位
12. Thread 与 Runnable
    Thread 是实现了 Runnable 接口的类

```JAVA
public interface Runnable {
    /**
     * When an object implementing interface <code>Runnable</code> is used
     * to create a thread, starting the thread causes the object's
     * <code>run</code> method to be called in that separately executing
     * thread.
     * <p>
     * The general contract of the method <code>run</code> is that it may
     * take any action whatsoever.
     *
     * @see     java.lang.Thread#run()
     */
    public abstract void run();
}

```

7. 如何处理线程返回值

- 主线程等待
- join 阻塞当前线程以等待子线程处理完毕
- Callable 接口实现:FutureTask 或线程池获取

8. 线程 6 状态

```JAVA
public enum State {
        /**
         * Thread state for a thread which has not yet started.
         */
        NEW,

        /**
         * Thread state for a runnable thread.  A thread in the runnable
         * state is executing in the Java virtual machine but it may
         * be waiting for other resources from the operating system
         * such as processor.
         */
        RUNNABLE,  // 包括操作系统Running和Ready  (R)

        /**
         * Thread state for a thread blocked waiting for a monitor lock.
         * A thread in the blocked state is waiting for a monitor lock
         * to enter a synchronized block/method or
         * reenter a synchronized block/method after calling
         * {@link Object#wait() Object.wait}.
         */
        BLOCKED,  // 等待获取排他锁

        /**
         * Thread state for a waiting thread.
         * A thread is in the waiting state due to calling one of the
         * following methods:
         * <ul>
         *   <li>{@link Object#wait() Object.wait} with no timeout</li>
         *   <li>{@link #join() Thread.join} with no timeout</li>
         *   <li>{@link LockSupport#park() LockSupport.park}</li>
         * </ul>
         *
         * <p>A thread in the waiting state is waiting for another thread to
         * perform a particular action.
         *
         * For example, a thread that has called {@code Object.wait()}
         * on an object is waiting for another thread to call
         * {@code Object.notify()} or {@code Object.notifyAll()} on
         * that object. A thread that has called {@code Thread.join()}
         * is waiting for a specified thread to terminate.
         */
        WAITING,  // 无限期等待，需要被显示唤醒

        /**
         * Thread state for a waiting thread with a specified waiting time.
         * A thread is in the timed waiting state due to calling one of
         * the following methods with a specified positive waiting time:
         * <ul>
         *   <li>{@link #sleep Thread.sleep}</li>
         *   <li>{@link Object#wait(long) Object.wait} with timeout</li>
         *   <li>{@link #join(long) Thread.join} with timeout</li>
         *   <li>{@link LockSupport#parkNanos LockSupport.parkNanos}</li>
         *   <li>{@link LockSupport#parkUntil LockSupport.parkUntil}</li>
         * </ul>
         */
        TIMED_WAITING,  // 限期等待，会被自动唤醒

        /**
         * Thread state for a terminated thread.
         * The thread has completed execution.
         */
        TERMINATED;   // 线程结束
    }
```

9. Thread.sleep 和 Object.wait 区别
   sleep 只会让出 CPU 不会导致锁行为改变 wait 会让出 CPU 并释放同步资源锁
   sleep 是 Thread 类的方法 wait 是 Object 类的方法
   sleep 到处用 wait 在 synchronized 里使用
10. Thread.yield 将线程状态从 running 变成 runnable
11. synchronized 锁的是对象
12. 异常处理机制:What(异常类型) Where(异常堆栈追踪) Why(异常信息)
    实现了 Throwable 接口
    Error:致命错误 **StackOverflowError**,**OutOfMemoryError** 等,无法处理
    Exception:程序可以处理的异常，可以处理
    - RuntimeException 不可预知的，例如数组越界 **IndexOutOfBoundsException**
    - 非 RuntimeExction 可预知 **IOException**
13. 异常处理机制:
    抛出异常=>
14. Collection 包括 List /Set /Queue
    Map 不是 Collection

```JAVA
public interface Collection<E> extends Iterable<E> {
  ...
}
```

15. Map(数组+链表+红黑树)
    java8 以后树化
    hashMap 获取 hash 到散列的过程

    注意：hashMap 不是线程安全的 可能两个线程同时 rehashing 造成死锁 且 rehashing 比较耗时

    ```JAVA
    static final float DEFAULT_LOAD_FACTOR = 0.75f;
    static final int DEFAULT_INITIAL_CAPACITY = 1 << 4; // aka 16
    static final int MAXIMUM_CAPACITY = 1 << 30;
    static final int TREEIFY_THRESHOLD = 8;
    static final int UNTREEIFY_THRESHOLD = 6;

    static final int hash(Object key) {
        int h;
        return (key == null) ? 0 : (h = key.hashCode()) ^ (h >>> 16);
    }
    ```

16. 线程安全的 hashMap
    synchronizedMap 内部有 mutex
    public 方法全部 synchronized
    串行执行 多线程时效率低

    ```JAVA
    private static class SynchronizedMap<K,V>
        implements Map<K,V>, Serializable {
        private static final long serialVersionUID = 1978198479659022715L;

        private final Map<K,V> m;     // Backing Map
        final Object      mutex;        // Object on which to synchronize

        SynchronizedMap(Map<K,V> m) {
            this.m = Objects.requireNonNull(m);
            mutex = this;
        }

        SynchronizedMap(Map<K,V> m, Object mutex) {
            this.m = m;
            this.mutex = mutex;
        }

        public int size() {
            synchronized (mutex) {return m.size();}
        }
        public boolean isEmpty() {
            synchronized (mutex) {return m.isEmpty();}
        }
        public boolean containsKey(Object key) {
            synchronized (mutex) {return m.containsKey(key);}
        }
        public boolean containsValue(Object value) {
            synchronized (mutex) {return m.containsValue(value);}
        }
        public V get(Object key) {
            synchronized (mutex) {return m.get(key);}
        }

        public V put(K key, V value) {
            synchronized (mutex) {return m.put(key, value);}
        }
        public V remove(Object key) {
            synchronized (mutex) {return m.remove(key);}
        }
        public void putAll(Map<? extends K, ? extends V> map) {
            synchronized (mutex) {m.putAll(map);}
        }
        public void clear() {
            synchronized (mutex) {m.clear();}
        }
    ...
    ```

17. CAS 比较并交换
18. ConcurrentHashMap
    ConcurrentHashMap 细粒度化了锁
    早期 ConcurrentHashMap:分段锁 Segment 来实现 数组+链表
    Java8 之后 CAS+synchronized 更加细化了锁 数组+链表+红黑树
19. JUC (concurrent 包)
20. BIO NIO AIO
    NonBlock-IO 多路复用的同步非阻塞 IO
21. Spring 框架
    Spring IOC 控制反转
    依赖注入:把底层类**作为参数传给上层类而不是在上层类使用 new 创建**，实现上层对下层的控制，也即创建与使用分离
    依赖倒置原则/IOC/**IOC 容器**/DO 的关系
    IOC 容器创建对象使用查 config 进行深度优先遍历

    BeanFactory 与 ApplicationContext

22. AOP
    关注点分离
