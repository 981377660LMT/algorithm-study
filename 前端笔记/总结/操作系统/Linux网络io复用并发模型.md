同一时刻，你只能被动的处理一个快递员的签收业务。其他快递员打电话打不进来，只能干瞪眼等待。那么解决这个问题，家里多买 N 个座机， 但是依然是你一个人接，也处理不过来，需要用影分身术创建都个自己来接电话(采用多线程或者多进程）来处理。

​ 这种方式就是没有多路 IO 复用的情况的解决方案， 但是在单线程计算机时代(无法影分身)，这简直是灾难。
https://www.kancloud.cn/aceld/golang/1958320

0. 流是什么： 流是数据的载体，可以进行 I/O 操作的内核对象
   文件、管道、套接字……
   流的入口：文件描述符(fd)
1. I/O 操作是什么：所有对流的读写操作，我们都可以称之为 IO 操作。当一个流中， 在没有数据 read 的时候，或者说在流中已经写满了数据，再 write，我们的 IO 操作就会出现一种现象，就是阻塞现象.
2. 如何解决阻塞死等待
   办法一：阻塞+多进程/多线程
   需要开辟线程浪费资源

   办法二：非阻塞、忙轮询

   ```Go
   while true {
     for i in 流[] {
       if i has 数据 {
         读 或者 其他处理
       }
     }
   }
   ```

   宏观上来看，是同时可以与多个快递员沟通(并发效果)、 但是快递员在于用户沟通时耽误前进的速度(**浪费 CPU**)。

   办法三：IO 多路复用：select，与平台无关
   我们可以开设一个代收网点，让快递员全部送到代收点。这个网店管理员叫 select。这样我们就可以在家休息了，麻烦的事交给 select 就好了。当有快递的时候，select 负责给我们打电话，期间在家休息睡觉就好了。
   **但 select 代收员比较懒**，她记不住快递员的单号，还有快递货物的数量。**她只会告诉你快递到了，但是是谁到的，你需要挨个快递员问一遍。**

   ```go
   while true {
     select(流[]); //阻塞

     //有消息抵达
     for i in 流[] {
       if i has 数据 {
         读 或者 其他处理
       }
     }
   }
   ```

   linux 只能监听 1024 个 fd

   办法四：IO 多路复用：epoll，仅限 linux
   epoll 的服务态度要比 select 好很多，在通知我们的时候，不仅告诉我们有几个快递到了，还分别告诉我们是谁谁谁。我们只需要按照 epoll 给的答复，来询问快递员取快递即可

   ```Go
   while true {
     可处理的流[] = epoll_wait(epoll_fd); //阻塞

     //有消息抵达，全部放在 “可处理的流[]”中
     for i in 可处理的流[] {
       读 或者 其他处理
     }
   }
   ```

   epoll？
   与 select，poll 一样，对 **I/O 多路复用**(同一个线程同时监听多个 IO)的技术
   只关心**活跃**的链接，无需遍历全部描述符集合
   能够处理大量的链接请求(**系统可以打开的文件数目**)

   > epoll 是解决 C10K 问题的利器，通过两个方面解决了 select/poll 的问题。

   epoll 在内核里使用「红黑树」来关注进程所有待检测的 Socket，红黑树是个高效的数据结构，增删查一般时间复杂度是 O(logn)，通过对这棵红黑树的管理，不需要像 select/poll 在每次操作时都传入整个 Socket 集合，减少了内核和用户空间大量的数据拷贝和内存分配。
   epoll 使用**事件驱动**的机制，内核里维护了一个「链表」来记录就绪事件，只将有事件发生的 Socket 集合传递给应用程序，不需要像 select/poll 那样轮询扫描整个集合（包含有和无事件的 Socket ），大大提高了检测的效率。

3. epoll 的 api

(1)创建 EPOLL

```C
/**
 * @param size 告诉内核监听的数目
 *
 * @returns 返回一个epoll句柄（即一个文件描述符）
 */
int epoll_create(int size);

int epfd = epoll_create(1000);
创建红黑树根节点epfd
```

(2) 控制 EPOLL

```C
/**
* @param epfd 用epoll_create所创建的epoll句柄
* @param op 表示对epoll监控描述符控制的动作
*
* EPOLL_CTL_ADD(注册新的fd到epfd)
* EPOLL_CTL_MOD(修改已经注册的fd的监听事件)
* EPOLL_CTL_DEL(epfd删除一个fd)
*
* @param fd 需要监听的文件描述符
* @param event 告诉内核需要监听的事件
*
* @returns 成功返回0，失败返回-1, errno查看错误信息
*/
int epoll_ctl(int epfd, int op, int fd,
struct epoll_event *event);
```

创建一个用户态的事件，绑定到某个 fd 上，然后添加到内核中的 epoll 红黑树中。

(3) 等待 EPOLL

```C
/**
*
* @param epfd 用epoll_create所创建的epoll句柄
* @param event 从内核得到的事件集合
* @param maxevents 告知内核这个events有多大,
* 注意: 值 不能大于创建epoll_create()时的size.
* @param timeout 超时时间
* -1: 永久阻塞
* 0: 立即返回，非阻塞
* >0: 指定微秒
*
* @returns 成功: 有多少文件描述符就绪,时间到时返回0
* 失败: -1, errno 查看错误
*/
int epoll_wait(int epfd, struct epoll_event *event,
							 int maxevents, int timeout);

使用
struct epoll_event my_event[1000];
int event_cnt = epoll_wait(epfd, my_event, 1000, -1);
```

epoll_wait 是一个阻塞的状态，如果内核检测到 IO 的读写响应，会抛给上层的 epoll_wait, **返回给用户态一个已经触发的事件队列**，同时阻塞返回。开发者可以从队列中取出事件来处理，其中事件里就有绑定的对应 fd 是哪个(之前添加 epoll 事件的时候已经绑定)。

(4) 使用 epoll 编程主流程骨架

```C
int epfd = epoll_crete(1000);

//将 listen_fd 添加进 epoll 中
epoll_ctl(epfd, EPOLL_CTL_ADD, listen_fd,&listen_event);

while (1) {
	//阻塞等待 epoll 中 的fd 触发
	int active_cnt = epoll_wait(epfd, events, 1000, -1);

	for (i = 0 ; i < active_cnt; i++) {
		if (evnets[i].data.fd == listen_fd) {
			//accept. 并且将新accept 的fd 加进epoll中.
		}
		else if (events[i].events & EPOLLIN) {
			//对此fd 进行读操作
		}
		else if (events[i].events & EPOLLOUT) {
			//对此fd 进行写操作
		}
	}
```

3. epoll 的两种触发模式 LT ET
   水平触发：
   水平触发的主要特点是，如果用户在监听 epoll 事件，当内核有事件的时候，会拷贝给用户态事件，但是如果用户只处理了一次，那么剩下没有处理的会在下一次 epoll_wait 再次返回该事件。
   **这样如果用户永远不处理这个事件，就导致每次都会有该事件从内核到用户的拷贝，耗费性能**，但是水平触发相对安全，最起码事件不会丢掉，除非用户处理完毕。
   边缘触发：
   边缘触发，相对跟水平触发相反，**当内核有事件到达， 只会通知用户一次，至于用户处理还是不处理，以后将不会再通知**。这样减少了拷贝过程，增加了性能，但是相对来说，如果用户马虎忘记处理，将会产生事件丢的情况。
   epoll 支持边缘触发和水平触发的方式，而 select/poll 只支持水平触发，一般而言，边缘触发的方式会比水平触发的效率高。
4. 单点 Server 的 7 种并发模型汇总
   https://www.kancloud.cn/aceld/golang/1958324

   1. 单线程 accept 无 io 复用：处理业务中，如果有新客户端 Connect 过来，Server 无响应，直到当前套接字全部业务处理完毕。

   2. 单线程 Accept+多线程读写业务（无 IO 复用）：创建链接成功，得到 Connfd1 套接字后，创建一个新线程 thread1 用来处理客户端的读写业务。main thead 依然回到 Accept 阻塞等待新客户端

   3. 单线程多路 IO 复用: 主线程 main thread 创建 listenFd 之后，采用**多路 I/O 复用机制(如:select、epoll)**进行 IO 状态阻塞监控。多路 I/O 复用阻塞，非忙询状态，不浪费 CPU 资源， CPU 利用率较高。**虽然可以监听多个客户端的读写状态，但是同一时间内，只能处理一个客户端的读写操作，实际上读写的业务并发为 1。**即：监测多个但是只能处理一个

   4. 单线程多路 IO 复用+多线程读写业务(业务工作池):main thread 按照固定的协议读取消息，并且交给 worker pool 工作线程池， 工作线程池在 server 启动之前就已经开启固定数量的 thread，里面的线程只处理消息业务，不进行套接字读写操作。**虽然多个 worker 线程处理业务，但是最后返回给客户端，依旧需要排队，因为出口还是 main thread 的 Read + Write**

   5. 单线程 IO 复用+多线程 IO 复用(链接线程池):**主流**
      类似 nodejs 的多进程主从模型
      main thread 采用多路 I/O 复用机制, **将新生成的 connFd1 分发给 Thread Pool 中的某个线程进行监听**。Thread Pool 中的每个 thread 都启动多路 I/O 复用机制(select、epoll),用来监听 main thread 建立成功并且分发下来的 socket 套接字。
      主线程只负责 accept 和传递 fd 给线程池
   6. 单进程 IO 复用+多进程 IO 复用(链接线程池)
      就是 nodejs/nginx 的多进程主从模型
      main process(主进程)**不再进行 Accept 操作**，而是将 Accept 过程分散到各个子进程(process)中.
      与现成的区别：关键是线程共享 fd 而**fd 只在当前进程可见**
      因为文件描述符实际上是进程内部打开文件表的下标
      **主进程需要分发 listenFd 让子进程去抢占处理 accept(listenFd)**
      main process 只是监听 ListenFd 状态，一旦触发读事件(有新连接请求). 通过一些 IPC(进程间通信：如信号、共享内存、管道)等, 让各自子进程 Process 竞争 Accept 完成链接建立，并各自监听。
      例子：
      Node.js 里通过 node app.js 开启一个服务进程，多进程就是进程的复制（fork），fork 出来的每个进程都拥有自己的独立空间地址、数据栈，一个进程无法访问另外一个进程里定义的变量、数据结构，只有建立了 IPC 通信，进程之间才可数据共享。
      nodejs 里进程之间可以借助**内置的 IPC 机制通信**
      父进程：
      接收事件 process.on('message')
      发送信息给子进程 master.send()
      子进程：
      接收事件 process.on('message')
      发送信息给父进程 process.send()
   7. 单线程多路 I/O 复用+多线程多路 I/O 复用+多线程
      Thread Pool 中的每个 thread 都启动多路 I/O 复用机制(select、epoll),用来监听 main thread 建立成功并且分发下来的 socket 套接字。一旦其中某个被监听的客户端套接字触发 I/O 读写事件,那么，**会立刻开辟一个新线程来处理 I/O 读写业务。**
      该模型过于理想化，因为要求 CPU 核心数量足够大。
      如果硬件 CPU 数量可数(目前的硬件情况)，那么该模型将造成大量的 CPU 切换成本浪费。因为为了保证读写并行通道与客户端 1:1 的关系，那么 Server 需要开辟的 Thread 数量就与客户端一致，那么线程池中做多路 I/O 复用的监听线程池绑定 CPU 数量将变得毫无意义。

## 有哪些常见的 IO 模型？

- 同步阻塞 IO（Blocking IO）：用户线程发起 IO 读/写操作之后，线程阻塞，直到可以开始处理数据；对 CPU 资源的利用率不够；
- 同步非阻塞 IO（Non-blocking IO）：发起 IO 请求之后可以立即返回，如果没有就绪的数据，需要不断地发起 IO 请求直到数据就绪；不断重复请求消耗了大量的 CPU 资源；
- IO 多路复用
- 异步 IO（Asynchronous IO）：用户线程发出 IO 请求之后，继续执行，由内核进行数据的读取并放在用户指定的缓冲区内，在 IO 完成之后通知用户线程直接使用。
