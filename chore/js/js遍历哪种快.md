https://leetcode.cn/problems/handling-sum-queries-after-update/submissions/

1. forEach (3468 ms)

```ts
function handleQuery(nums1, nums2, queries) {
  let sum = nums2.reduce((a, b) => a + b, 0)
  const res = []
  queries.forEach(([op, a, b]) => {
    if (op === 1) {
      for (let i = a; i <= b; i++) {
        nums1[i] ^= 1
      }
    } else if (op === 2) {
      sum += a * nums1.reduce((a, b) => a + b, 0)
    } else {
      res.push(sum)
    }
  })
  return res
}
```

2. for of (5000 ms)

```ts
function handleQuery(nums1, nums2, queries) {
  let sum = nums2.reduce((a, b) => a + b, 0)
  const res = []
  for (const [op, a, b] of queries) {
    if (op === 1) {
      for (let i = a; i <= b; i++) {
        nums1[i] = 1 - nums1[i]
      }
    } else if (op === 2) {
      sum += a * nums1.reduce((a, b) => a + b, 0)
    } else {
      res.push(sum)
    }
  }
  return res
}
```

3. for i + 多次取值 (5000 ms)

```ts
function handleQuery(nums1, nums2, queries) {
  let sum = nums2.reduce((a, b) => a + b, 0)
  const res = []
  for (let i = 0; i < queries.length; i++) {
    const [op, a, b] = queries[i]
    if (op === 1) {
      for (let i = a; i <= b; i++) {
        nums1[i] ^= 1
      }
    } else if (op === 2) {
      sum += a * nums1.reduce((a, b) => a + b, 0)
    } else {
      res.push(sum)
    }
  }
  return res
}
```

所以不需要返回和 break 时,forEach 最快
一般来说,朴素的 for 循环最快，但是这里涉及到多个下标获取和解构,forEach 更快
