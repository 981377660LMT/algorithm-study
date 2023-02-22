https://leetcode.cn/problems/handling-sum-queries-after-update/submissions/

1. forEach (3468 ms)

```ts
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
```

2. for of (4772 ms)

```ts
 for (const [op, a, b] of queries) {
    if (op === 1) {
      for (let i = a; i <= b; i++) {
        nums1[i] = 1-nums1[i]
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
for (let i = 0; i < queries.length; i++) {
  const op = queries[i][0]
  const a = queries[i][1]
  const b = queries[i][2]
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
```

所以不需要返回和 break 时,forEach 最快
