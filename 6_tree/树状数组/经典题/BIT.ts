/**
 * @summary
 * 高效计算数列的前缀和，区间和
 * 树状数组或二叉索引树（Binary Indexed Tree, Fenwick Tree）
 * 性质
 * 1. tree[x]保存以x为根的子树中叶节点值的和
 * 2. tree[x]的父节点为tree[x+lowbit(x)]
 * 3. tree[x]节点覆盖的长度等于lowbit(x)
 * 4. 树的高度为logn+1
 */
class BIT1 {
  readonly size: number
  private readonly tree = new Map<number, number>()

  constructor(size: number) {
    this.size = size
  }

  add(index: number, delta: number): void {
    if (index <= 0) throw Error('index 应为正整数')
    for (let i = index; i <= this.size; i += this.lowbit(i)) {
      this.tree.set(i, (this.tree.get(i) ?? 0) + delta)
    }
  }

  query(index: number): number {
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i -= this.lowbit(i)) {
      res += this.tree.get(i) ?? 0
    }
    return res
  }

  sumRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }

  private lowbit(x: number): number {
    return x & -x
  }
}

class BIT2 {
  readonly size: number
  private readonly tree1 = new Map<number, number>()
  private readonly tree2 = new Map<number, number>()

  constructor(size: number) {
    this.size = size
  }

  add(left: number, right: number, k: number): void {
    this._add(left, k)
    this._add(right + 1, -k)
  }

  query(left: number, right: number): number {
    return this._query(right) - this._query(left - 1)
  }

  private _add(index: number, delta: number): void {
    if (index <= 0) throw Error('查询索引应为正整数')
    for (let i = index; i <= this.size; i += this.lowbit(i)) {
      this.tree1.set(i, (this.tree1.get(i) ?? 0) + delta) // 此处进行了差分操作，记录差分操作大小
      this.tree2.set(i, (this.tree2.get(i) ?? 0) + (index - 1) * delta) // 前x-1个数没有进行差分操作，这里把总值记录下来
    }
  }

  private _query(index: number): number {
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i -= this.lowbit(i)) {
      res += index * (this.tree1.get(i) ?? 0) - (this.tree2.get(i) ?? 0)
    }

    return res
  }

  private lowbit(x: number): number {
    return x & -x
  }
}

if (require.main === module) {
  const bit1 = new BIT1(5)
  console.log(bit1.query(1)) // 0
  bit1.add(1, 3)
  console.log(bit1.query(1)) // 3

  const bit2 = new BIT2(10)

  bit2.add(2, 4, 1) // 区间更新
  bit2.add(2, 2, 1) // 单点更新
  console.log(bit2)
  console.log(bit2.query(2, 4)) // 区间查询
  console.log(bit2.query(2, 2)) // 单点查询

  function maximumWhiteTiles(tiles: number[][], carpetLen: number): number {
    const bit = new BIT2(Math.max(...tiles.flat()) + 10)
    for (const [left, right] of tiles) bit.add(left, right, 1)
    let res = 0
    for (const [left] of tiles) res = Math.max(res, bit.query(left, left + carpetLen - 1))
    return res
  }
}

export { BIT1, BIT2 }
