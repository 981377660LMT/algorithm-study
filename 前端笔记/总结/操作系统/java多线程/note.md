http://www.imooc.com/wiki/concurrencylesson/yield.html
https://docs.qq.com/doc/DSVNyZ2FNWWFkeFpO 面试题
https://docs.qq.com/doc/DSVp5R1NubFpERVZ1 思维导图

1. 并发的实质是一个物理 CPU 在若干道程序之间多路复用，其目的是提高有限物理资源的运行效率。
2. 在多线程编程实践中，线程的个数往往多于 CPU 的个数，所以一般都称多线程并发编程而不是多线程并行编程。
3. 2 种开启线程的方法：
   **实现 Runnable 接口**/继承 Thread

```JAVA
Thread 类中
@Override
public void run() {
    if (target != null) {
        target.run();
    }
}
target 是 Runnable
```

第一种是重写 run 方法，第二种实现 Runnable 接口的 run 方法，然后再把该 runnable 实例传给 Thread 类。除此之外，从表面上看线程池、定时器等工具类也可以创建线程，但是它们的本质都逃不出刚才所说的范围。

4. 实现 Runnable 接口和继承 Thread 类哪种方式更好?

- runnable 因为具体的任务 run 应该和创建解耦;且 java 单继承

6. 停止线程的方式
   interrupt
7. Thread 与 Object 重要方法
   sleep 在 Thread
   wait notify notifyAll 在 Object:
   wait 释放锁
8. 多线程 join/yield 方法
   多线程环境下，如果需要确保**某一线程执行完毕后**才可继续执行后续的代码，就可以通过使用 join 方法完成这一需求设计。
9. 操作系统是为每个线程分配一个时间片来占有 CPU
   正常情况下当一个线程把分配给自己的时间片使用完后，线程调度器才会进行下一轮的线程调度，这里所说的 “自己占有的时间片” 即 CPU 分配给线程的执行权。
   当一个线程通过某种可行的方式向操作系统提出让出 CPU 执行权时，就是在告诉线程调度器自己占有的时间片中**还没有使用完的部分自己不想使用了**，**主动放弃剩余的时间片**，并在合适的情况下，重新获取新的执行时间片。

```JAVA
public static native void yield();
一个 Native Method 就是一个 Java 调用的非 Java 代码的接口
该方法的实现由非 java 语言实现。
可以理解为操作调用操作系统的方法接口。
```

5. 线程属性
   ID Name(清晰有意义的名字) isDaemon Priority
6. yield 方法和 sleep 方法的区别

   1. sleep () 方法给其他线程运行机会时不考虑线程的优先级，因此会给低优先级的线程以运行的机会；
      yield () 方法只会给相同优先级或更高优先级的线程以运行的机会；
   2. 线程执行 sleep () 方法后转入 ( WAITING ) 状态，而执行 yield () 方法后转入就绪 (ready) 状态；
   3. sleep () 方法声明会抛出 InterruptedException, 而 yield () 方法没有声明任何异常；

7. 为了让用户感觉多个线程是在同时执行的， CPU 资源的分配采用了**时间片轮转**的策略，也就是给每个线程分配一个时间片，线程在时间片内占用 CPU 执行任务。
   什么是线程的上下文切换：当前线程使用完时间片后，就会处于就绪状态并让出 CPU，让其他线程占用，这就是上下文切换，从当前线程的上下文切换到了其他线程。在切换线程上下文时需要保存当前线程的执行现场， 当再次执行时根据保存的执行现场信息恢复执行现场。
   **线程上下文切换时机**：当前线程的 CPU 时间片使用完或者是当前线程被其他线程中断时
8. Java 中的线程分为两类，分别为 daemon 线程（守护线程〉和 user 线程（用户线程）。
   在 JVM 启动时会调用 main 函数， main 函数所在的线程就是一个用户线程，其实在 JVM 内部同时还启动了**好多守护线程，比如垃圾回收线程**。
   只要任何非守护线程还在运行，程序就不会终止。
   只要有一个用户线程还没结束， 正常情况下 JVM 就不会退出。
   Linux 守护进程是系统级别的，当系统退出时，才会终止。
   而 Java 中的守护线程是 JVM 级别的，当 JVM 中无任何用户进程时，守护进程销毁，JVM 退出，程序终止。
   创建方式：将线程转换为守护线程可以通过调用 Thread 对象的 **setDaemon (true)** 方法来实现。
   **守护线程应该永远不去访问固有资源，如文件、数据库，因为它会在任何时候甚至在一个操作的中间发生中断。**

```JAVA
public class DemoTest {
    public static void main(String[] args) throws InterruptedException {
        Thread threadOne = new Thread(new Runnable() {
            @Override
            public void run() {
                //代码执行逻辑
            }
        });
        threadOne.setDaemon(true); //设置threadOne为守护线程
        threadOne. start();
    }
}

```

    守护线程的作用及使用场景
    正常开发过程中，一般心跳监听，垃圾回收，临时数据清理等通用服务会选择守护线程。

8. ThreadLocal **线程隔离**
   ThreadLocal 很容易让人望文生义，想当然地认为是一个 “本地线程”。其实，ThreadLocal 并不是一个 Thread，**而是 Thread 的局部变量**，也许把它命名为 ThreadLocalVariable 更容易让人理解一些。
   **当多个线程操作这个变量时，实际操作的是自己本地内存里面的变量**，从而避免了线程安全问题。
   ThreadLocal 是线程本地存储，**在每个线程中都创建了一个 ThreadLocalMap 对象**，**每个线程可以访问自己内部 ThreadLocalMap 对象内的 value**。通过这种方式，避免资源在多线程间共享。
   多线程下的 ThreadLocal 才是它存在的真实意义
   在很多情况下，ThreadLocal 比直接使用 synchronized 同步机制解决线程安全问题更简单，更方便，且结果程序拥有更高的并发性。
   flask 的线程隔离机制
9. 多线程 AVO 原则
   A：即 Atomic，原子性操作原则。对基本数据类型的变量读和写是保证原子性的，要么都成功，要么都失败，这些操作不可中断。

   V：即 volatile，可见性原则。后续的小节会对 volatile 关键字进行深入的讲解，此处只需要理解概念性问题即可。使用 volatile 关键字，保证了变量的可见性，到主存拿数据，不是到缓存里拿。

   O：即 order， 就是有序性。代码的执行顺序，在代码编译前的和代码编译后的执行顺序不变。

10. 抛开语言，谈操作系统的线程的生命周期及线程 5 种状态，这是我们学习 Java 多线程 6 种状态的基础；
    我们来看下 Java 线程的 6 种状态的概念。

新建 (New)：实现 Runnable 接口或者继承 Thead 类可以得到一个线程类，new 一个实例出来，线程就进入了初始状态。

运行 (Running)：线程调度程序从可运行池中选择一个线程作为当前线程时线程所处的状态。这也是线程进入运行状态的唯一方式。(**包括 操作系统的 Running 和 Runnable**)

阻塞 (Blocked)：**阻塞状态是线程在进入 synchronized** 关键字修饰的方法或者代码块时，由于其他线程正在执行，不能够进入方法或者代码块而被阻塞的一种状态。

等待 (Waiting)：执行 wait () 方法后线程进入等待状态，如果没有显示的 notify () 方法或者 notifyAll () 方法唤醒，该线程会一直处于等待状态。

超时等待 (Timed_Waiting)：执行 sleep（Long time）方法后，线程进入超时等待状态，时间一到，自动唤醒线程。

终止状态 (Terminal)：当线程的 run () 方法完成时，或者主线程的 main () 方法完成时，我们就认为它终止了。这个线程对象也许是活的，但是，它已经不是一个单独执行的线程。线程一旦终止了，就不能复生。

11. synchronized 关键字
    synchronized 同步块是 Java 提供的**一种原子性内置锁**，Java 中的每个对象都可以把它当作一个同步锁来使用，这些 Java 内置的使用者看不到的锁被称为内部锁，也叫作监视器锁。
    代码在进入 synchronized 代码块前会自动获取内部锁，这时候其他线程访问该同步代码块时会被阻塞挂起。拿到内部锁的线程会在正常退出同步代码块或者抛出异常后或者在同步块内调用了该内置锁资源的 wait 系列方法时释放该内置锁。
12. 在实现生产者消费者问题时，可以采用三种方式：

    - 使用 Object 的 wait/notify 的消息通知机制，本节课程我们采用该方式结合 synchronized 关键字进行生产者与消费者模式的实现；
    - 使用 Lock 的 Condition 的 await/signal 的消息通知机制；
    - 使用 BlockingQueue 实现。本文主要将这三种实现方式进行总结归纳。

13. volatile 关键字
    当其他线程读取该共享变量时，会从主内存重新获取最新值，而不是使用当前线程的工作内存中的值。
14. CAS
    synchronized 时代效率问题:获得锁的线程在运行，其他被挂起的线程只能等待着
    volatile 时代原子操作问题: volatile 不会造成阻塞,但是 volatile 不能保证原子性
    CAS 操作诞生的意义：CAS（Compare And Swap 比较和交换）**解决了 volatile 不能保证原子性的问题**。从而 CAS 操作即能够解决锁的效率问题，也能够保证操作的原子性。(java.util.concurrent (JUC java 并发工具包) 就是建立在 CAS 之上的)
    CAS 主要包含三个操作数，内存位置 V，进行比较的原值 A，和新值 B。
    **当位置 V 的值与 A 相等时，CAS 才会通过原子方式用新值 B 来更新 V**，否则不会进行任何操作。无论位置 V 的值是否等于 A，都将返回 V 原有的值。
