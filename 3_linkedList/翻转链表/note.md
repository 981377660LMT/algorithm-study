基本 api 反转链表**reverse**

K 个一组反转就是

1. 初始化 dummyP p1 p2
2. 找到每个 K 长度的首位 p1 p2
3. 断开这一段(p2.next=undefined) 并且 反转这一段**reverse**(p1)
4. 虚拟头接上去 dummyP.next=p2
5. dummyP=p1 继续下一次反转
