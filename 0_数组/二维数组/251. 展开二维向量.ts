class Vector2D {
  private inner: number
  private outer: number
  private vector: number[][]

  constructor(vec: number[][]) {
    this.inner = 0
    this.outer = 0
    this.vector = vec
  }

  next(): number {
    this.mightUpdate()
    const res = this.vector[this.outer][this.inner]
    this.inner++
    return res
  }

  hasNext(): boolean {
    this.mightUpdate()
    return this.outer < this.vector.length
  }

  // 到顶了就跳到下一个数组(如果为空 直接跳过)
  private mightUpdate() {
    while (this.outer < this.vector.length && this.inner === this.vector[this.outer].length) {
      this.outer++
      this.inner = 0
    }
  }
}

/**
 * Your Vector2D object will be instantiated and called as such:
 * var obj = new Vector2D(vec)
 * var param_1 = obj.next()
 * var param_2 = obj.hasNext()
 */
//  请记得 重置 在 Vector2D 中声明的类变量（静态变量），
//  因为类变量会 在多个测试用例中保持不变，影响判题准确

// 迭代器的主要目的之一就是最小化辅助空间的使用。
// 我们应尽可能的利用现有的数据结构，
// 只需要添加足够多的额外空间来跟踪下一个值。
// 在某些情况下，我们要遍历的数据结构太大，甚至无法放入内存（例如文件系统）
