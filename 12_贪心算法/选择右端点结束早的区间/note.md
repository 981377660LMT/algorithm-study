基本思路是贪心选择结束早的区间
重叠区间问题模板
**按右端点为第一关键字，长度为第二关键字排序，然后从前往后遍历线段**，优先选择结尾最短的区间,留给后面的空间更多

```TS
// 找到最长的不重叠区间
var findLongestChain = function (pairs: number[][]) {
  // sort by the earliest finish time
  pairs.sort((a, b) => a[1] - b[1])
  let prev = pairs[0],
    chain = 1

  for (let i = 1; i < pairs.length; i++) {
    const [prevS, prevE] = prev
    const [currS, currE] = pairs[i]
    if (prevE < currS) {
      prev = pairs[i]
      chain++
    }
  }
  return chain
}
```
