# 解析 Golang 网络 IO 模型之 EPOLL

![alt text](image.png)

## 1. IO 多路复用

1. 何为 IO 多路复用
   多路复用指的是，由一个执行单元，同时对多个对象提供服务，形成一对多的服务关系
   ![alt text](image-1.png)
   • 多路：存在多个需要处理 io event 的 fd（linux 中，一切皆文件，所有事物均可抽象为一个文件句柄 file descriptor，简称 fd）
   • 复用：复用一个 loop thread ，同时为多个 fd 提供处理服务（线程 thread 是内核视角下的最小调度单位；多路复用通常为循环模型 loop model，因此称为 loop thread）
2. IO 多路复用的简单实现
