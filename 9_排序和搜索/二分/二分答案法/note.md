力扣上，对于某些运行效率比较高的语言，不想算二分上界的话，可以直接设成 `1e16`，因为 JS 的 MAX_SAFE_INTEGER 2^53-1 < 1e16

---

这种解法不容易被想到
问题特征是：一般是求第 k 个 或者是直接求一个数的值
并且：`求确定值很难，求范围很容易`

**数据范围很大** 考虑二分答案  
`1802. 有界数组中指定下标处的最大值.py ` 范围是
`1 <= n <= maxSum <= 10^9`

check 函数 O(n) =>双指针
`668. 乘法表中第k小的数-双指针.ts`
`719找出第K小的距离对-双指针.ts`

**找第 k 小=>计数最左二分**

```Python
        def countNGT(mid) -> int:
            """"目标值`小于等于`mid的答案数"""
            res = 0
            for right in range(len(nums)):
             ...
            return res

        left, right = 0, int(1e18)
        while left <= right:
            mid = (left + right) >> 1
            if count(mid) < k:
                left = mid + 1
            # 大于等于k right都左移
            else:
                right = mid - 1
        return left

        python 3.10以后 等价于
        bisect_left(range(int(1e18)), k, key=countNGT)
```

第 k 小 第 k 大的数 简单中等题很多是优先队列(多路归并)，难题很多是二分答案
