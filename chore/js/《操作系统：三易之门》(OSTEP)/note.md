# 《操作系统：三易之门》(OSTEP)

https://github.com/remzi-arpacidusseau/ostep-translations/tree/master/chinese
TODO: 南京大学蒋老师课程

本书围绕 3 个主题元素展开讲解：虚拟化（virtualization）、并发（concurrency）和持久性（persistence）。

操作系统实际上做了什么：
它取得 CPU、内存或磁盘等物理资源（resources），甚对它们进行虚拟化（virtualize）。
它处理与并发（concurrency）有关的麻烦且棘手的问题。
它持久地（persistently）存储文件，从而使它们长期安全。
鉴于我们希望建立这样一个系统，所以要有一些目标，以帮助我们集中设计和实现，并在必要时进行折中。
找到合适的折中是建立系统的关键。

---

## intro

1. Dialogue
2. Introduction

## virtualization (虚拟化)

以最基本的计算机资源 CPU 为例，假设一个计算机只有一个 CPU（尽管现代计算机一般拥有 2 个、4 个或者更多 CPU），虚拟化要做的就是将这个 CPU 虚拟成多个虚拟 CPU 并分给每一个进程使用。
因此，`每个应用都以为自己在独占 CPU，但实际上只有一个 CPU。这样操作系统就创造了美丽的假象——它虚拟化了 CPU。`

3. Dialogue

4. Processes（进程）

时分共享（time sharing）CPU 技术，允许用户如愿运行多个并发进程。潜在的开销就是性能损失，
因为如果 CPU 必须共享，每个进程的运行就会慢一点。
通过将代码和静态数据加载到内存中，通过创建和初始化栈以及执行与 I/O 设置相关的
其他工作，OS 现在（终于）为程序执行搭好了舞台。然后它有最后一项任务：启动程序，
在入口处运行，即 main()。

5. Process API (进程 API)

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
