模板

1. pq push 起始值
2. 获取 k 个结果
3. push 下一个

```JS
 pq.push(起始值)
 while (pq.length && res.length < k) {
    const [_, i, j] = pq.shift()!
    res.push([nums1[i], nums2[j]])
    // ...
    pushNext(val,i,j+1)
  }
 return res
```
