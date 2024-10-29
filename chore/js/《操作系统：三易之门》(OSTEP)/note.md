# 《操作系统：三易之门》(OSTEP)

Operating Systems: Three Easy Pieces
https://book.douban.com/subject/19973015/

https://github.com/remzi-arpacidusseau/ostep-translations/tree/master/chinese
https://pages.cs.wisc.edu/~remzi/OSTEP/Chinese/
TODO: 南京大学蒋老师课程

本书围绕 3 个主题元素展开讲解：虚拟化（virtualization）、并发（concurrency）和持久性（persistence）。

操作系统实际上做了什么：
它取得 CPU、内存或磁盘等物理资源（resources），甚对它们进行虚拟化（virtualize）。
它处理与并发（concurrency）有关的麻烦且棘手的问题。
它持久地（persistently）存储文件，从而使它们长期安全。
鉴于我们希望建立这样一个系统，所以要有一些目标，以帮助我们集中设计和实现，并在必要时进行折中。
找到合适的折中是建立系统的关键。

TODO：问题答案

---

## intro

1. Dialogue
2. Introduction(操作系统介绍)
   https://pages.cs.wisc.edu/~remzi/OSTEP/Chinese/toc.pdf
   - 虚拟化：
     操作系统是对物理资源的代理。提供的api叫system call。
     - 虚拟化CPU => 系统拥有非常多的虚拟 CPU 的假象
     - 虚拟化内存 => 每个进程访问自己的`私有虚拟地址空间`（virtual address space）（有时称为地址空间，address space），操作系统以某种方式映射到机器的物理内存上。`就好像每个正在运行的程序都有自己的私有内存，而不是与其他正在运行的程序共享相同的物理内存`
   - 并发
   - 持久性
     操作系统中操理磁盘的软件通常称为`文件系统（file system）`。因此它负责以可靠和高效的方式，将用户创建的任何文件（file）存储在系统的磁盘上。
     调用操作系统的 open、write、close。
     大多数文件系统都包含某种复杂的写入协议，以确保在系统崩溃时不会丢失数据。
   - 设计目标
     操作系统实际上做了什么：它取得 CPU、内存或磁盘等物理资源
     （resources），甚对它们进行虚拟化（virtualize）。它处理与甚发（concurrency）有关的麻烦且棘手的问题。它持久地（persistently）存储文件，从而使它们长期安全。
     - 系统设计是权衡的艺术：抽象vs性能
     - 令一个目标是在应用程序之间以及在 OS 和应用程序之间提供保护（protection）。让进程彼此隔离是保护的关键。
     - 操作系统往往力求提供高度的可靠性（reliability）。
   - 简单历史
     - 早期操作系统：只是一些库
     - 超越库：保护。添加一些特殊的硬件指令和硬件状态，让向操作系统过渡变为`更正式的、受控的过程。`
     - 多道程序时代
       操作系统不是一次只运行一项作业，而是将大量作业加载到内存中甚在它们之间快速切换，从而提高 CPU 利用率
     - 摩登时代
       UNIX 的重要性

## virtualization (虚拟化)

以最基本的计算机资源 CPU 为例，假设一个计算机只有一个 CPU（尽管现代计算机一般拥有 2 个、4 个或者更多 CPU），虚拟化要做的就是将这个 CPU 虚拟成多个虚拟 CPU 并分给每一个进程使用。
因此，`每个应用都以为自己在独占 CPU，但实际上只有一个 CPU。这样操作系统就创造了美丽的假象——它虚拟化了 CPU。`

3. Dialogue

4. Processes（进程）

操作系统的最基本抽象：进程
进程就是运行中的程序。操作系统为正在运行的程序提供的抽象。
进程的状态：内存(地址空间)+寄存器+I/O

`时分共享（time sharing）`CPU 技术，允许用户如愿运行多个并发进程。潜在的开销就是性能损失，因为如果 CPU 必须共享，每个进程的运行就会慢一点。

- 进程api分类

  - create
  - destroy
  - wait
  - miscellaneous control (其他控制)
  - statu

- 进程创建
  通过将代码和静态数据加载到内存中，通过创建和初始化栈以及执行与 I/O 设置相关的其他工作，OS 现在（终于）为程序执行搭好了舞台。然后它有最后一项任务：启动程序，
  在入口处运行，即 main()。
  ![alt text](image.png)
- 进程状态
  决策由操作系统调度程序完成
  ![alt text](image-1.png)
  - Running
  - Ready
  - Blocked
  - Initial
  - Final (僵尸状态，已退出但尚未清理)
- 数据结构
  进程列表保存三种状态的进程

5. Process API (进程 API)

UNIX 系统采用了一种非常有趣的创建新进程的方式，即通过一对系统调用：`fork()和 exec()`。进程还可以通过第三个系统调用 `wait()`，来等待其创建的子进程执行完成。

- fork
  创建新进程
  注意：fork 会复制父进程的内存，包括代码、数据和堆栈。`父进程获得的返回值是子进程的 PID，而子进程获得的返回值是 0。`
  CPU 调度程序（scheduler）决定了某个谁刻哪个进程被执行
- wait
  有时父进程需要等待子进程执行完毕, 用wait或者waitpid。
- exec
  **子进程执行与父进程不同的程序**
  给我可执行程序的名称（如 wc）及需要的参数（如 p3.c）后，exec()会从可执行程序中加载代码和静态数据
- 为什么这样设计api
  为什么设计如此奇怪的接口，来完成简单的、创建新进程的任务?
  这种`分离 fork()及 exec()的做法`在构建 UNIX shell 的
  时候谁非常有用，因为这给了 shell 在 `fork 之后 exec 之前运行代码的机会`，这些代码可以谁运行新程序前改变环境，从而让一系列有趣的功能很容易实现.
  相当于before钩子。
  例如：shell实现重定向命令，完成`子进程创建后，调用exec之前，shell关闭了标准输出`，打开了文件xxx.txt。这样子进程的输出就会写入到xxx.txt文件中。

  管道命令。shell创建两个子进程，一个执行ls，一个执行grep，然后通过管道连接两个子进程。

- 其他api
  kill 向进程发送信号
  ps 显示进程列表
  top 显示进程列表和资源使用情况

6. Direct Execution (直接执行)

7. CPU Scheduling (CPU 调度)

8. Multi-level Feedback (多级反馈)

9. Lottery Scheduling (抽奖调度)

10. Multi-CPU Scheduling (多 CPU 调度)

11. Summary

12. Dialogue

13. Address Spaces (地址空间)

14. Memory API (内存 API)

15. Address Translation (地址转换)

16. Segmentation (分段)

17. Free Space Management (空闲空间管理)

18. Introduction to Paging (分页简介)

19. Translation Lookaside Buffers (TLB)

20. Advanced Page Tables (高级页表)

21. Swapping: Mechanisms (交换：机制)

22. Swapping: Policies (交换：策略)

23. Complete VM Systems (完整的虚拟机系统)

24. Summary

## concurrency (并发)

25. Dialogue

26. Concurrency and Threads (并发和线程)

27. Thread API (线程 API)

28. Locks (锁)

29. Locked Data Structures (锁定数据结构)

30. Condition Variables (条件变量)

31. Semaphores (信号量)

32. Concurrency Bugs (并发 Bug)

33. Event-based Concurrency (基于事件的并发)

34. Summary

## persistence (持久性)

35. Dialogue

36. I/O Devices (I/O 设备)

37. Hard Disk Drives (硬盘驱动器)

38. Redundant Disk Arrays (RAID)

39. Files and Directories (文件和目录)

40. File System Implementation (文件系统实现)

41. Fast File System (FFS) (快速文件系统)

42. FSCK and Journaling (文件系统检查和日志)

43. Log-Structured File System (LFS) (日志结构文件系统)

44. Data Integrity and Protection (数据完整性和保护)

45. Summary

46. Dialogue

47. Distributed Systems (分布式系统)

48. Network File System (NFS) (网络文件系统)

49. Andrew File System (AFS) (安德鲁文件系统)

50. Summary

## appendices

- Virtual Machines

- Monitors

- Lab Tutorial

- Systems Labs

- xv6 Labs
