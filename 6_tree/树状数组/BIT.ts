interface IBIT {
  add: (x: number, k: number) => void
  query: (x: number) => number
}

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
class BinaryIndexedTree implements IBIT {
  private nums: number[]
  private tree: number[]
  private size: number

  constructor(nums: number[]) {
    this.nums = nums
    this.size = nums.length
    this.tree = Array(this.size + 1).fill(0)
    // tree的0号位置不存值;初始化「树状数组」，要默认数组是从 1 开始
    for (let i = 0; i < this.size; i++) {
      this.add(i + 1, nums[i])
    }
  }

  /**
   *
   * @param x
   * @param k
   * @description
   * 单点修改：tree数组下标x处的值加k
   */
  add(x: number, k: number) {
    if (x === 0) return
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree[i] += k
    }
  }

  /**
   *
   * @param x
   * @description
   * 区间查询：返回前x项的和
   */
  query(x: number) {
    let res = 0
    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += this.tree[i]
    }
    return res
  }

  /**
   *
   * @param left  nums下标
   * @param right  nums下标
   * @returns 区域和
   */
  sumRange(left: number, right: number): number {
    return this.query(right + 1) - this.query(left)
  }

  /**
   *
   * @param x
   * @returns x 的二进制表示中，最低位的 1 的位置。
   * @example
   * ```js
   * console.log(3 & -3) // 1
   * console.log(4 & -4) // 4
   *
   * ```
   */
  private lowbit(x: number) {
    return x & -x
  }
}

if (require.main === module) {
  const bit = new BinaryIndexedTree([1, 2, 3, 4, 5])
  console.log(bit.query(1))
  bit.add(1, 3)
  console.log(bit.query(1))
}

export { BinaryIndexedTree }

// 利用数组实现前缀和，查询本来是 O(1)，但是对于频繁更新的数组，每次重新计算前缀和，时间复杂度 O(n)。
// 此时树状数组的优势便立即显现。
