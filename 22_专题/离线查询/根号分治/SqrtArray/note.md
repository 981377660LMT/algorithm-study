参考分块的 SortedList
还可以实现区间删除/翻转等操作

---

js 的特点是

1. 操作指针(对象引用)特别慢
2. 遍历/操作数组/splice 特别快 (和 golang 差不多)
3. 空间占用不能太大

这种特点，使得 js 非常适合根号分块的数据结构(相比较平衡树等数据结构而言)

类似的思想
https://leetcode.cn/problems/subrectangle-queries/solution/bu-bao-li-geng-xin-ju-zhen-yuan-su-de-jie-fa-by-li/
