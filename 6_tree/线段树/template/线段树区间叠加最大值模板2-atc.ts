/**
 * @see {@link https://atcoder.jp/contests/typical90/submissions/23947318}
 */
abstract class AbstractSegmentTree<T = number> {
  /**
   * @description {@link tree} 的初始值
   */
  protected abstract initTreeValue(): T

  /**
   *@description 左右结点的合并方式
   */
  protected abstract mergeChildren(v1: T, v2: T): T

  /**
   * @description 数组值的更新方式
   */
  protected abstract updateTreeValue(oldV: T, newV: T): T

  private readonly upper: number
  private readonly tree: T[]
  private offset: number

  constructor(lenOrArray: number | T[]) {
    this.upper = typeof lenOrArray === 'number' ? lenOrArray : lenOrArray.length
    this.offset = 1
    while (this.offset < this.upper) {
      this.offset *= 2
    }

    const bufsize = this.offset * 2
    this.tree = Array(bufsize)
    this.tree[0] = this.initTreeValue()

    if (Array.isArray(lenOrArray)) {
      for (let i = 0; i < this.upper; i++) {
        this.tree[this.offset + i] = lenOrArray[i]
      }

      for (let i = this.upper; i < this.offset; i++) {
        this.tree[this.offset + i] = this.initTreeValue()
      }

      for (let i = this.offset - 1; i >= 1; i--) {
        this.tree[i] = this.mergeChildren(this.tree[i * 2], this.tree[i * 2 + 1])
      }
    } else {
      for (let i = 1; i < bufsize; i++) {
        this.tree[i] = this.initTreeValue()
      }
    }
  }

  /**
   * @param index 数组下标 0 <= index < n
   * @param value 更新的值 更新方式取决于 {@link updateTreeValue}
   */
  update(index: number, value: T): void {
    let k = index + this.offset
    this.tree[k] = this.updateTreeValue(this.tree[k], value)
    while (k > 1) {
      k >>= 1
      this.tree[k] = this.mergeChildren(this.tree[2 * k], this.tree[2 * k + 1])
    }
  }

  /**
   * @returns  切片`[left:right]`内的信息 0 <= left <= right <=n
   */
  query(left: number, right: number): T {
    let res: T = this.initTreeValue()
    for (left += this.offset, right += this.offset; left < right; left >>= 1, right >>= 1) {
      if (left & 1) {
        res = this.mergeChildren(res, this.tree[left++])
      }

      if (right & 1) {
        res = this.mergeChildren(res, this.tree[--right])
      }
    }

    return res
  }

  get length() {
    return this.upper
  }
}

/**
 * @description 线段树RMQ最大值(快速版)
 */
class MaxSegmentTree2 extends AbstractSegmentTree<number> {
  protected initTreeValue(): number {
    return -Infinity
  }

  protected mergeChildren(a: number, b: number): number {
    return Math.max(a, b)
  }

  protected updateTreeValue(oldV: number, newV: number): number {
    return newV + oldV
    return newV
  }
}

/**
 * @description 线段树RMQ最小值(快速版)
 */
class MinSegmentTree2 extends AbstractSegmentTree<number> {
  protected initTreeValue(): number {
    return Infinity
  }

  protected mergeChildren(a: number, b: number): number {
    return Math.min(a, b)
  }

  protected updateTreeValue(oldV: number, newV: number): number {
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
