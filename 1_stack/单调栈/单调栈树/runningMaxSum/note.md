# RunningMaxSum

- 定义
  `runningMax[i] = max( nums[L], ..., nums[i] ) (L <= i <= R)`
  `runningMaxSum = sum( runningMax[L..R] ).`
- 练习
  [3420. 统计 K 次操作以内得到非递减子数组的数目](https://leetcode.cn/problems/count-non-decreasing-subarrays-after-k-operations/)

  - [x] 单调栈树+树上倍增
  - [ ] 滑窗 + 单调栈
