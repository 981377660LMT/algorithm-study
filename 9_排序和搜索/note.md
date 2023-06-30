快速排序是二叉树的前序遍历，归并排序是二叉树的后序遍历

```C++
void sort(int[] nums, int lo, int hi) {
    /****** 前序遍历位置 ******/
    // 通过交换元素构建分界点 p
    int p = partition(nums, lo, hi);
    /************************/

    sort(nums, lo, p - 1);
    sort(nums, p + 1, hi);
}
```

```C++
void sort(int[] nums, int lo, int hi) {
    int mid = (lo + hi) / 2;
    sort(nums, lo, mid);
    sort(nums, mid + 1, hi);

    /****** 后序遍历位置 ******/
    // 合并两个排好序的子数组
    merge(nums, lo, mid, hi);
    /************************/
}
```

JS sort 目前都是稳定的排序

1. JS sort 从 ES10（EcmaScript 2019）开始，要求 Array.prototype.sort 为稳定排序。
2. 除 IE 外所有浏览器都已支持稳定的 sort 排序
3. ES6 规范要求 sort 实现原理：https://tc39.es/ecma262/#sec-array.prototype.sort
4. `chrome 70 开始使用的排序算法为 Timsort（python sorted 的排序算法内部实现也基于 Timsort）`
5. Timsort 是稳定的排序，实现原理简单来说 ：
   按照升序(后一元素大于等于前一元素，a[i] <= a[i+1])和严格降序(前一元素大于后一元素 a[i]>a[i+1])的规则，将原来的数组分解为若干个 run，升序的 run 就保持不变，严格降序的 run 就翻转，最终得到若干个升序的 run。
   合并相邻的 run，直到只剩下一个排序好的 run。细节：
   合并的时候使用倍增搜寻法+二分查找法
   维护了一个 stack,这个栈会记录起始的索引位置和每个 run 的长度，满足后入栈的 run 长度大于前两个入栈的 run 的长度之和，让 run 的长度递减，避免长度差太多的 run 合并。
   时间复杂度：
   平均情况 O(n log n)
   最坏情况 O(n log n)
   最好情况 O(n)
   空间复杂度：O(n)

---

https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/_note/heuristic_search.md#L1-L17

# Heuristic Search Algorithms

## Overview

Basically, there are three choices:

- A\*
- IDA* (iterative deepening A*)
- RBFS (recursive best first search)

If the state space is `sufficiently small`, use A*.
Otherwise, use IDA* or RBFS.
If good solutions are spreaded among search pathes, use IDA\*.
Otherwise, i.e., good solutions are condensed, use RBFS.

---

golang 的 sort 源码值得学习

https://segmentfault.com/a/1190000039668324
https://www.huawei.com/cn/open-source/blogs/optimizelab-overlay-developer-guide
