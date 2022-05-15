// https://leetcode-cn.com/problems/fancy-sequence/solution/ru-guo-bu-liao-jie-cheng-fa-ni-yuan-jiu-hao-hao-xu/

import assert from 'assert'

const MOD = BigInt(1e9 + 7)
// 注意1e9+7*(1e9+7) 超出了js的精度了 9e15
// 结点实现超内存了 换数组

class SegmentTreeNode {
  left = -1n
  right = -1n
  isLazy = false
  lazyAdd = 0n
  lazyMul = 1n
  value = 0n
}

class SegmentTree {
  private tree: SegmentTreeNode[]

  constructor(size: number) {
    this.tree = Array.from({ length: size << 2 }, () => new SegmentTreeNode())
    this.build(1, 1, size)
  }

  update(root: number, left: number, right: number, delta: bigint, type: 'ADD' | 'MUL'): void {
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      if (type === 'ADD') {
        node.isLazy = true
        node.value += (node.right - node.left + 1n) * delta
        node.value %= MOD
        node.lazyAdd += delta
        node.lazyAdd %= MOD
      } else {
        node.isLazy = true
        node.value *= delta
        node.value %= MOD
        node.lazyAdd *= delta
        node.lazyAdd %= MOD
        node.lazyMul *= delta
        node.lazyMul %= MOD
      }

      return
    }

    this.pushDown(root)
    const mid = (node.left + node.right) >> 1n
    if (left <= mid) this.update(root << 1, left, right, delta, type)
    if (mid < right) this.update((root << 1) | 1, left, right, delta, type)
    this.pushUp(root)
  }

  query(root: number, left: number, right: number): number {
    const node = this.tree[root]
    if (left <= node.left && node.right <= right) {
      return Number(node.value % MOD)
    }

    this.pushDown(root)
    const mid = (node.left + node.right) >> 1n

    // 单点查询

    if (left <= mid) return this.query(root << 1, left, right)
    else return this.query((root << 1) | 1, left, right)
  }

  private build(root: number, left: number, right: number): void {
    const node = this.tree[root]
    node.left = BigInt(left)
    node.right = BigInt(right)

    if (left === right) {
      return
    }

    const mid = Number((node.left + node.right) >> 1n)
    this.build(root << 1, left, mid)
    this.build((root << 1) | 1, mid + 1, right)
    // this.pushUp(root) // 这里不用
  }

  /**
   * @param root 向下传递懒标记和懒更新的值 `isLazy`, `lazyValue`，并用 `lazyValue` 更新子区间的值
   */
  private pushDown(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]

    if (node.isLazy) {
      left.isLazy = true
      left.value = left.value * node.lazyMul + node.lazyAdd
      left.lazyAdd = left.lazyAdd * node.lazyMul + node.lazyAdd
      left.lazyMul *= node.lazyMul

      right.isLazy = true
      right.value = right.value * node.lazyMul + node.lazyAdd
      right.lazyAdd = right.lazyAdd * node.lazyMul + node.lazyAdd
      right.lazyMul *= node.lazyMul

      left.value %= MOD
      left.lazyAdd %= MOD
      left.lazyMul %= MOD
      right.value %= MOD
      right.lazyAdd %= MOD
      right.lazyMul %= MOD

      node.isLazy = false
      node.lazyAdd = 0n
      node.lazyMul = 1n
    }
  }

  /**
   * @description 只有单点查询 没必要pushUp
   */
  private pushUp(root: number): void {
    // const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    // node.value = left.value + right.value
    // node.value %= MOD
  }
}

class Fancy {
  private size: number
  private tree: SegmentTree

  /**
   * 总共最多会有 1e5 次对 append，addAll，multAll 和 getIndex 的调用。
   */
  constructor() {
    this.size = 0
    this.tree = new SegmentTree(1e5 + 10)
  }

  /**
   * @param val 将整数 val 添加在序列末尾
   */
  append(val: number): void {
    this.size += 1
    this.tree.update(1, this.size, this.size, BigInt(val), 'ADD')
  }

  /**
   * @param inc 将所有序列中的现有数值都增加 inc
   */
  addAll(inc: number): void {
    if (this.size === 0) return
    this.tree.update(1, 1, this.size, BigInt(inc), 'ADD')
  }

  /**
   * @param m 将序列中的所有现有数值都乘以整数 m
   */
  multAll(m: number): void {
    if (this.size === 0) return
    this.tree.update(1, 1, this.size, BigInt(m), 'MUL')
  }

  /**
   * @param idx 得到下标为 idx 处的数值,并将结果对 109 + 7 取余。
   * 如果下标大于等于序列的长度，请返回 -1
   */
  getIndex(idx: number): number {
    if (idx >= this.size) return -1
    return this.tree.query(1, idx + 1, idx + 1)
  }
}

/**
 * Your Fancy object will be instantiated and called as such:
 * var obj = new Fancy()
 * obj.append(val)
 * obj.addAll(inc)
 * obj.multAll(m)
 * var param_4 = obj.getIndex(idx)
 */

if (require.main === module) {
  const fancy = new Fancy()
  fancy.append(2)
  fancy.addAll(3)
  assert.strictEqual(fancy.getIndex(0), 5)

  fancy.append(7)
  assert.strictEqual(fancy.getIndex(1), 7) // 0

  fancy.multAll(2)
  assert.strictEqual(fancy.getIndex(0), 10)
  assert.strictEqual(fancy.getIndex(1), 14)
}

export {}
