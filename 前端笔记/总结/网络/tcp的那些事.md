https://coolshell.cn/articles/11564.html

1. TCP 协议的定义
   TCP 在网络 OSI 的七层模型中的第四层——Transport 层，IP 在第三层——Network 层，ARP 在第二层——Data Link 层，在第二层上的数据，我们叫 Frame，在第三层上的数据叫 Packet，第四层的数据叫 Segment。我们程序的数据首先会打到 TCP 的 Segment 中，然后 TCP 的 Segment 会打到 IP 的 Packet 中，然后再打到以太网 Ethernet 的 Frame 中，传到对端后，各个层解析自己的协议，然后把数据交给更高层的协议处理。
   四个非常重要的东西
   **Sequence Number** 是包的序号，用来解决网络包乱序（reordering）问题。
   **Acknowledgement Number** 就是 ACK——用于确认收到，用来解决不丢包的问题。
   **Window** 又叫 Advertised-Window，也就是著名的滑动窗口（Sliding Window），用于解决流控的。
   **TCP Flag** ，也就是包的类型，主要是用于操控 TCP 的状态机的。
2. 丢包时的重传机制
3. TCP 的流迭
4. TCP 滑动窗口
5. 拥塞处理
   拥塞控制主要是四个算法：1）慢启动，2）拥塞避免，3）拥塞发生，4）快速恢复
