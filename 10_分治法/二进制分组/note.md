# 二进制分组(Binary Grouping)

https://zerol.me/2017/08/04/group-by-binary/
https://rsk0315.hatenablog.com/entry/2019/06/19/124528

题目要求支持修改和查询操作。但是你只有一种`不支持修改`的算法，往往是通过预处理来支持快速查询操作。

- 把需要维护的元素按照个数的二进制表示，分成不超过 log𝑛 个组，构建每个组时进行一些预处理，使得每个组内都可以快速地询问信息。如果增加一个操作，分组的变化通过暴力删除和重建进行（当然如果可以合并的话那就只需要合并了）。(如 5 = 4 + 1, 6 = 4 + 2，从 5 增加到 6 就把原来那个大小为 1 的组删了，再将原来最后一个和新增加的那一个建一个大小为 2 的组)

- 与其他分治算法相比，**优势在于可以在线**
- 在线版 cdq 分治：以一个 log 的代价，让一个需要支持动态修改的问题变成不需要支持动态修改的问题。

https://www.cnblogs.com/Dfkuaid-210/p/bit_divide.html

https://www.cnblogs.com/TianMeng-hyl/p/14989441.html

[二进制分组的本质](https://www.mina.moe/archives/12681)

[尝试用二进制分组实现某些数据结构及习题](https://zhuanlan.zhihu.com/p/35519230)
