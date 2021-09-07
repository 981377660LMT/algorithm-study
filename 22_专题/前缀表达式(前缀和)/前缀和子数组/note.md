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
