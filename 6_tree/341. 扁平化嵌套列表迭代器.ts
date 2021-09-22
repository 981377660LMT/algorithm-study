interface NestedInteger {
  isInteger: () => boolean
  getInteger: () => number | null
  setInteger: (value: number) => unknown
  add: (elem: NestedInteger) => unknown
  getList: () => NestedInteger[]
}

class NestedIterator {
  private list: NestedInteger[]
  private index: number
  constructor(nestedList: NestedInteger[]) {
    this.list = []
    this.index = 0
    this.flatten(nestedList)
  }

  hasNext(): boolean {
    return this.index < this.list.length
  }

  next(): number {
    return this.list[this.index++].getInteger()!
  }

  // 核心
  private flatten(nestedList: NestedInteger[]) {
    for (const nestedInteger of nestedList) {
      if (nestedInteger.isInteger()) {
        this.list.push(nestedInteger)
      } else {
        this.flatten(nestedInteger.getList())
      }
    }
  }
}

// 输入：nestedList = [[1,1],2,[1,1]]
// 输出：[1,1,2,1,1]
// 解释：通过重复调用 next 直到 hasNext 返回 false，next 返回的元素的顺序应该是: [1,1,2,1,1]。
