dp 的话会超时
需要用二分

LIS 模板

- dp 法(都能求,O(n^2))

```JS
// dp[i]：到nums[i]为止的最长递增子序列长度
const dp = Array(nums.length).fill(1)

  for (let i = 1; i < nums.length; i++) {
    // 状态转移方程
    for (let j = 0; j < i; j++) {
      if (nums[j] < nums[i]) dp[i] = Math.max(dp[i], dp[j] + 1)
    }
  }
```

- 二分法(只能求 LIS 长度,O(nlogn))

```TS
const lengthOfLIS = function (nums: number[]): number {
  if (nums.length <= 1) return nums.length
  const LIS: number[] = [nums[0]]

  for (let i = 1; i < nums.length; i++) {
    if (nums[i] > LIS[LIS.length - 1]) {
      LIS.push(nums[i])
    } else {
      LIS[bisectLeft(LIS, nums[i])] = nums[i]
    }
  }
  console.log(LIS)
  return LIS.length
}

// 辅助函数
const bisectLeft = (arr: number[], target: number) => {
  let l = 0
  let r = arr.length - 1
  while (l <= r) {
    const mid = (r + l) >> 1
    if (arr[mid] === target) {
      r--
    } else if (arr[mid] > target) {
      r = mid - 1
    } else {
      l = mid + 1
    }
  }

  return l
}
```

1713. 得到子序列的最少操作次数
      当一个数组元素各不相同，LCS 问题可以转换为 LIS 问题
      当其中一个数组元素各不相同时，这时候每一个“公共子序列”都对应一个不重复元素数组的下标数组“上升子序列”，反之亦然
      O(n^2) 变 O(nlogn)

1714. 1751. 2008. 带权区间选择问题

`模板参考 1751`

```
1.按结束时间排序
2.每个区间为结尾初始化dp
3.bisectRight找到之前的最右区间，根据 存在/不存在pre 列转移方程

```

严格递增
插入时 bisectLeft
`1_最长上升子序列dp加二分解法.ts`
不严格递增
插入时 bisectRight
`1964. 找出到每个位置为止最长的有效障碍赛跑路线.py`

---

O(nloglogn)求 LIS

https://rsk0315.hatenablog.com/entry/2023/01/22/032133
