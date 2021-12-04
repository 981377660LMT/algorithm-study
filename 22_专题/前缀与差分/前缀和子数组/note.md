<!-- 需要考虑初始位置的0 -->
<!-- 保证[0,k]的正确性 即第一个元素从头开始 -->

```JS
  模板

  const pre = new Map<number, number>([[0, -1]])  // 存储前缀和和需要解决的量
  let sum = 0
  let res = 0

  for (let i = 0; i < nums.length; i++) {
    sum += nums[i]
    ...
    if (pre.has(sum)) ...
    else pre.set(sum, i)
  }

  return res
```

1. 子数组和为目标值，子数组和被某数整除:记录`模 k 的值 -> index 或者是 count`
2. 最长的子数组:记录前缀和**第一次出现的位置**（setdefault） 之后相等时跟哈希表里的对比
   `1124. 表现良好的最长时间段.py`
   `325. 和等于 k 的最长子数组长度.js`
   `525. 连续数组.ts`

3. 满足 i<j 的条件
   遍历中更新 pre，记录最新的位置
