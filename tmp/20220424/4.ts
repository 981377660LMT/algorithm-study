class BIT {
  public size: number
  private tree1: Map<number, number>
  private tree2: Map<number, number>

  constructor(size: number) {
    this.size = size
    // 记录∑di，即记录差分操作的大小
    this.tree1 = new Map()
    this.tree2 = new Map()
  }

  // 上下车问题
  addRange(left: number, right: number, k: number): void {
    this.add(left, k)
    this.add(right + 1, -k)
  }

  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }

  add(index: number, delta: number): void {
    if (index <= 0) throw Error('查询索引应为正整数')
    for (let i = index; i <= this.size; i += this.lowbit(i)) {
      this.tree1.set(i, (this.tree1.get(i) ?? 0) + delta) // 此处进行了差分操作，记录差分操作大小
      this.tree2.set(i, (this.tree2.get(i) ?? 0) + (index - 1) * delta) // 前x-1个数没有进行差分操作，这里把总值记录下来
    }
  }

  // 差分数组(上下车问题):
  // 假设现在有一个原数组a(假设a[0] = 0)，有一个数组d，d[i] = a[i] - a[i-1]，那么
  // a[i] = d[1] + d[2] + .... + d[i]
  // 差分数组 diff[i]，存储的是 res[i] - res[i - 1]；而差分数组 diff[0...i] 的和(树状数组更新/查询就是做求和)，就是 res[i] 的值。
  // 则[1,x]范围的和：a[1] + a[2] + a[3] + ... + a[x] = d[1] + d[1] + d[2] + d[1] + d[2] + d[3] + ... + d[1] + d[2] + d[3] + ... + d[x]
  // =x*(Σd[i]) - ∑(i-1)*d[i] (i从1到x)
  query(index: number): number {
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i -= this.lowbit(i)) {
      res += index * (this.tree1.get(i) ?? 0) - (this.tree2.get(i) ?? 0)
    }

    return res
  }

  private lowbit(x: number) {
    return x & -x
  }
}

function fullBloomFlowers(flowers: number[][], persons: number[]): number[] {
  const bit = new BIT(1e9 + 7)
  for (const [x, y] of flowers) {
    bit.addRange(x, y, 1)
  }

  const res = []
  for (const p of persons) {
    res.push(bit.queryRange(p, p))
  }

  return res
}
