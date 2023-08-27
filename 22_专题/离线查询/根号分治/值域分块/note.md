值域分块，顾名思义就是在值域上分块，似乎一般用来均摊复杂度，
比如有些时候会以 `O(1)修改 O(sqrt(n))查询` 代替两者均为 O(logn)的权值树状数组。

一般配合莫队使用，因为莫队中修改进行 O(nsqrt(n))次，但是查询只进行 O(n)次，并不均匀，因为会用值域分块代替树状数组，会用**修改复杂度低 O(1)的数据结构。**

类似 SortedList，区别在于**值域分块可以 O(1)定位到 pos 和 index，不需要二分查找。**

https://www.luogu.com.cn/problem/P4867
https://www.luogu.com.cn/problem/P4119
https://livinfly.top/186&decompose_9problems_for_beginner
https://www.cnblogs.com/flashhu/p/8437062.html
https://loj.ac/p?keyword=%E5%88%86%E5%9D%97%E5%85%A5%E9%97%A8
https://www.cnblogs.com/chifan-duck/p/17060540.html

---

[浅谈根号算法](https://ddosvoid.github.io/2020/10/18/%E6%B5%85%E8%B0%88%E6%A0%B9%E5%8F%B7%E7%AE%97%E6%B3%95/)
[莫队套值域分块](https://www.cnblogs.com/zaza-zt/p/15041167.html)
主要用途是维护一堆数，支持 O(1)插入、删除，以及 O(√N)复杂度实现查询前驱、后继、k 小、x 的排名（类似于平衡树）。

- 动态区间第 k 小
