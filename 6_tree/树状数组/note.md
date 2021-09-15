区间问题：
简单求区间和，用「前缀和」
多次将某个区间变成同一个数，用「线段树」
其他情况，用「树状数组」

```TS
class BIT {
  private size: number
  private tree: number[]

  constructor(size: number) {
    this.size = size
    this.tree = Array(size + 1).fill(0)
  }

  // 最好x都离散化正数
  add(x: number, k: number) {
    if (x <= 0) throw Error('查询索引应为正整数')
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree[i] += k
    }
  }

  query(x: number) {
    let res = 0
    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += this.tree[i]
    }
    return res
  }

  sumRange(left: number, right: number) {
    return this.query(right + 1) - this.query(left)
  }

  private lowbit(x: number) {
    return x & -x
  }
}
```

# 注意

一般来说
处理的数组需要进行离散化

```JS
  const set = new Set(nums)
  const map = new Map<number, number>()
  for (const [key, realValue] of [...set].sort((a, b) => a - b).entries()) {
    map.set(realValue, key + 1)  // key+1是因为查询和修改的树状数组的索引应为正整数
  }
  // Map(4) { 1 => 1, 2 => 2, 5 => 3, 6 => 4 }
```
