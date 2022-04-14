这种解法不容易被想到
问题特征是：一般是求第 k 个 或者是直接求一个数的值
并且：`求确定值很难，求范围很容易`

**数据范围很大** 考虑二分答案  
`1802. 有界数组中指定下标处的最大值.py ` 范围是
`1 <= n <= maxSum <= 10^9`

check 函数 O(n) =>双指针
`668. 乘法表中第k小的数-双指针.ts`
`719找出第K小的距离对-双指针.ts`

**找第 k 小模板=>计数最左二分**

```Python
        def count(mid) -> int:
            """"目标值`小于等于`mid的答案数"""
            res = 0
            for right in range(len(nums)):
             ...
            return res

        left, right = 1, int(1e20) # 注意左边不要0，因为有些地方除以mid会出错
        while left <= right:
            mid = (left + right) >> 1
            # 小于k时，移动left
            if count(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left
```
