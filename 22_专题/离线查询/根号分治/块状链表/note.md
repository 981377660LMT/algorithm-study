给一个序列，初始时为空，要求支持两种操作：

在某个位置之前插入一个数；
求某个区间中的数异或 x 的最大值。

# 块状链表(区块链)

块状链表一般在用平衡树实现代码非常复杂时使用

https://oi-wiki.org/ds/block-list/

块状链表，顾名思义，就是把分块和链表结合起来的神奇数据结构。
分块区间操作复杂度优秀，但是不能支持 插入/删除 操作。
链表单点插入删除复杂度优秀，但是不能支持大规模的区间操作。
但是两者相结合，就会变得非常厉害。

https://www.cnblogs.com/LcyRegister/p/17047026.html
https://www.luogu.com/article/jk8hy8ti 【模板】普通平衡树
https://www.cnblogs.com/chenxiaoran666/p/Luogu4278.html 【洛谷 4278】带插入区间 K 小值（块状链表+值域分块）
P4008 [NOI2003]文本编辑器

---

https://www.acwing.com/blog/content/28060/
