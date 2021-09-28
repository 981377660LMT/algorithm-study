class ZigzagIterator {
  private curIndex: number
  private list: number[]

  // 有点像归并排序的merge过程
  constructor(v1: number[], v2: number[]) {
    this.curIndex = 0
    this.list = []

    let i = 0
    while (i < v1.length && i < v2.length) {
      this.list.push(v1[i])
      this.list.push(v2[i])
      i++
    }

    while (i < v1.length) {
      this.list.push(v1[i])
      i++
    }

    while (i < v2.length) {
      this.list.push(v2[i])
      i++
    }
  }

  // 交替返回它们中间的元素
  next(): number {
    return this.list[this.curIndex++]
  }

  hasNext(): boolean {
    return this.curIndex !== this.list.length
  }
}

export {}
// 拓展：假如给你 k 个一维向量呢？
// 你的代码在这种情况下的扩展性又会如何呢?
// 使用n个数组
