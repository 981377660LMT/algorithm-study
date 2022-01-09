// 区间修改区间查询模板

// 树状数组天生用来动态维护数组前缀和
// 树状数组区间更新：我们想对某个区间进行更新，我们只需对其差分数组的区间两端进行更新，并在查询时进行前缀求和即可
// 通过维护两个树状数组来实现
class BIT {
  public size: number
  private tree1: number[]
  private tree2: number[]

  constructor(size: number) {
    this.size = size
    // 记录∑di，即记录差分操作的大小
    this.tree1 = Array(size + 1).fill(0)
    // 记录∑(i-1)*di，即记录这个位置之前没有进行差分操作的个数*差分操作大小
    // 树状数组add就是对每个1位加入差分信息
    // query就是对每个1位查询这些信息的影响总和
    this.tree2 = Array(size + 1).fill(0)
    // a1+a2+...+ar=r*∑di-∑(i-1)*di
  }

  // 上下车问题
  addRange(left: number, right: number, k: number): void {
    this.add(left, k)
    this.add(right + 1, -k)
  }

  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }

  private add(x: number, k: number): void {
    if (x <= 0) throw Error('查询索引应为正整数')
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree1[i] += k // 此处进行了差分操作，记录差分操作大小
      this.tree2[i] += (x - 1) * k // 前x-1个数没有进行差分操作，这里把总值记录下来
    }
  }

  // 差分数组(上下车问题):
  // 假设现在有一个原数组a(假设a[0] = 0)，有一个数组d，d[i] = a[i] - a[i-1]，那么
  // a[i] = d[1] + d[2] + .... + d[i]
  // 差分数组 diff[i]，存储的是 res[i] - res[i - 1]；而差分数组 diff[0...i] 的和(树状数组更新/查询就是做求和)，就是 res[i] 的值。
  // 则[1,x]范围的和：a[1] + a[2] + a[3] + ... + a[x] = d[1] + d[1] + d[2] + d[1] + d[2] + d[3] + ... + d[1] + d[2] + d[3] + ... + d[x]
  // =x*(Σd[i]) - ∑(i-1)*d[i] (i从1到x)
  private query(x: number): number {
    let res = 0

    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += x * this.tree1[i] - this.tree2[i]
    }

    return res
  }

  private lowbit(x: number) {
    return x & -x
  }
}

if (require.main === module) {
  const bit = new BIT(10)
  // console.log(bit.query(2))
  // bit.add(2, 1)
  bit.addRange(2, 4, 1) // 区间更新
  bit.addRange(2, 2, 1) // 单点更新
  console.log(bit)
  console.log(bit.queryRange(2, 4)) // 区间查询
  console.log(bit.queryRange(2, 2)) // 单点查询
  // console.log(bit.query(3))
  // console.log(bit.query(4))
  // console.log(bit.query(5))
}

export { BIT }
