树状数组支持区间修改的前提是`运算可逆`，`有逆运算才能差分`。
加法有逆运算减法，乘法有逆运算除法，异或的逆运算是自身，这些都可以用。像最大值运算，位与运算，矩阵乘法之类就用不了树状数组了，还是要用线段树。
**只要是半群信息都可以用线段树维护。而树状数组则要求信息可以差分，ST 表要求是幂等半群信息。**

# 注意

一般来说需要将用于查询/修改的所有值需要进行离散化(set+并排序，map 映射成树状数组的索引，相对大小不变)

```JS
  const set = new Set(nums)
  const map = new Map<number, number>()
  for (const [key, realValue] of [...set].sort((a, b) => a - b).entries()) {
    map.set(realValue, key + 1)  // key+1是因为查询和修改的树状数组的索引应为正整数
  }
  // Map(4) { 1 => 1, 2 => 2, 5 => 3, 6 => 4 }
```

附上:逆序对问题三种解法

1. 手动维护一个有序的数组(java-treeset,python-sortedList) 是 O(n^2)还是 O(nlogn)取决于使用的数据结构**插入**的复杂度
2. 归并排序的性质
3. 树状数组

tree[i] = a[i-lowbit(i)+1] + ... + a[i]
