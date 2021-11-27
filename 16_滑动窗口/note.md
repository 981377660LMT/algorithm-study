滑动窗口的滑动需要具有单向性:例如移动右一直增，移动左一直减

思路都是用 map 记录数字对应的索引/出现次数

变长滑窗模板：

```JS
  let res = 0
  let sum = 0
  let left = 0

  for (let right = 0; right < array.length; right++) {
    // 处理新元素
    sum+=...

    // 需要调整
    while (sum>k) {
      ...
      left++
    }

    // 结算
    res=Math.max(res,right-left+1)
  }

  return res
```
