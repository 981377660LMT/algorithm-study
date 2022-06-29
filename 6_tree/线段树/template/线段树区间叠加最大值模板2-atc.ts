/**
 * @see {@link https://atcoder.jp/contests/typical90/submissions/23947318}
 */
abstract class AbstractSegmentTree<T = number> {
  /**
   * @description {@link tree} 的初始值
   */
  protected abstract getDefaultValue(): T

  /**
   *@description 左右结点的合并方式
   */
  protected abstract merger(v1: T, v2: T): T

  /**
   * @description 数组值的更新方式
   */
  protected abstract updater(oldV: T, newV: T): T

  private readonly len: number
  private readonly tree: T[]
  private n: number

  constructor(lenOrArray: number | T[]) {
    this.len = typeof lenOrArray === 'number' ? lenOrArray : lenOrArray.length
    this.n = 1
    while (this.n < this.len) {
      this.n *= 2
    }

    const bufsize = this.n * 2
    this.tree = Array(bufsize)
    this.tree[0] = this.getDefaultValue()

    if (Array.isArray(lenOrArray)) {
      for (let i = 0; i < this.len; i++) this.tree[this.n + i] = lenOrArray[i]
      for (let i = this.len; i < this.n; i++) this.tree[this.n + i] = this.getDefaultValue()
      for (let i = this.n - 1; i >= 1; i--)
        this.tree[i] = this.merger(this.tree[i * 2], this.tree[i * 2 + 1])
    } else {
      for (let i = 1; i < bufsize; i++) this.tree[i] = this.getDefaultValue()
    }
  }

  get length() {
    return this.len
  }

  /**
   * @param index 数组下标 0 <= index < n
   * @param value 更新的值 更新方式取决于 {@link updater}
   */
  update(index: number, value: T): void {
    let k = index + this.n
    this.tree[k] = this.updater(this.tree[k], value)
    while (k > 1) {
      k >>= 1
      this.tree[k] = this.merger(this.tree[2 * k], this.tree[2 * k + 1])
    }
  }

  /**
   * @returns  切片`[left:right]`内的信息 0 <= left <= right <=n
   */
  query(left: number, right: number): T {
    let res: T = this.getDefaultValue()
    for (left += this.n, right += this.n; left < right; left >>= 1, right >>= 1) {
      if (left & 1) {
        res = this.merger(res, this.tree[left++])
      }

      if (right & 1) {
        res = this.merger(res, this.tree[--right])
      }
    }

    return res
  }
}

/**
 * @description 线段树RMQ最大值(快速版)
 */
class MaxSegmentTree2 extends AbstractSegmentTree<number> {
  protected getDefaultValue(): number {
    return -Infinity
  }

  protected merger(a: number, b: number): number {
    return Math.max(a, b)
  }

  protected updater(oldV: number, newV: number): number {
    return newV + oldV
    return newV
  }
}

/**
 * @description 线段树RMQ最小值(快速版)
 */
class MinSegmentTree2 extends AbstractSegmentTree<number> {
  protected getDefaultValue(): number {
    return Infinity
  }

  protected merger(a: number, b: number): number {
    return Math.min(a, b)
  }

  protected updater(oldV: number, newV: number): number {
    return newV + oldV
    return newV
  }
}

if (require.main === module) {
  const tree = new MaxSegmentTree2([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  console.log(tree.length)
  console.log(tree.query(0, 1))
  tree.update(0, 3)
  console.log(tree.query(0, 1))
}

export { MaxSegmentTree2, MinSegmentTree2 }
