1. 远程连接命令 ssh scp
   ssh user@ip
2. 查看本地网络状态
   ifconfig:状态信息包含 IP 地址、子网掩码、MAC 地址等
   netstat:查看目前本机的网络使用情况
   查看端口占用
   netstat -ntlp | grep 80
3. 网络测试
   ping
   telnet:想知道对方服务器是否有对应该端口的进程
4. DNS 查询
   host
   dig
5. HTTP
   curl

最后，列一下本次提到的 Linux 下常用的网络命令：

远程登录的 ssh 指令；

远程传输文件的 scp 指令；

查看网络接口的 ifconfig 指令；

查看网络状态的 netstat 指令；

测试网络延迟的 ping 指令；

可以交互式调试和服务端的 telnet 指令；

两个 DNS 查询指令 host 和 dig；

可以发送各种请求包括 HTTPS 的 curl 指令。
