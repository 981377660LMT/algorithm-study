# 将第一个元素移动到数组末尾

对于每一个询问项 query，我们想求出它在排列 P 中的位置，实际上只要知道 `它的前面有几个数` 就可以了。
在现实生活中，你要管理一个队伍，想把其中的一个人放到队首，你会怎么做？最简单的做法是直接把**这个人拉出来让他/她站到第一个人的前面就行了**。为什么我们可以这么做？这是因为现实生活中的队伍不是数组，**在第一个人前面是有很多空间的**。
我们知道查询的次数 Q，那么我们可以**使用一个长度为 Q+M 的数组**，一开始的排列 P 占据了数组的最后 M 个位置，而每处理一个询问项 query，我们将其直接放到数组的前 Q 个位置就行了，顺序是从右往左放置。

```
对于排列 [1, 2, 3, 4, 5] 以及查询 [3, 1, 2, 1]，一开始的数组为：

_ _ _ _ 1 2 3 4 5
前面空出了四个位置，即查询的长度。

我们第一次查询 3，3 之前有 2 个数。随后将 3 移到前面：

_ _ _ 3 1 2 _ 4 5
我们第二次查询 1，1 之前有 1 个数。随后将 1 移到前面：

_ _ 1 3 _ 2 _ 4 5
我们第二次查询 2，1 之前有 2 个数。随后将 2 移到前面：

_ 2 1 3 _ _ _ 4 5
我们第二次查询 1，1 之前有 1 个数。随后将 1 移到前面：

1 2 _ 3 _ _ _ 4 5
```
