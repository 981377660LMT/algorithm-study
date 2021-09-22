// 拉链法
// 定长拉链数组
class MyHashMap1 {
  private bucketsNum: number
  private itemsPerBucket: number
  private table: number[][]

  constructor() {
    this.bucketsNum = 1000
    this.itemsPerBucket = 1001 // 质数
    this.table = Array.from({ length: this.itemsPerBucket }, () => [])
  }

  put(key: number, value: number): void {
    const col = this.getHash(key)
    const row = this.getPosition(key)
    this.table[row][col] = value
  }

  get(key: number): number {
    const col = this.getHash(key)
    const row = this.getPosition(key)
    return this.table[row][col] === undefined ? -1 : this.table[row][col]
  }

  remove(key: number): void {
    const col = this.getHash(key)
    const row = this.getPosition(key)
    this.table[row][col] = -1
  }

  // 矩阵的列
  private getHash(key: number) {
    return key % this.bucketsNum
  }

  // 矩阵的行
  private getPosition(key: number) {
    return ~~(key / this.bucketsNum)
  }
}

// 不定长拉链数组
class MyHashMap2 {
  private bucketsNum: number
  private table: [number, number][][]

  constructor() {
    this.bucketsNum = 1009
    this.table = Array.from<number, [number, number][]>({ length: this.bucketsNum }, () => [])
  }

  put(key: number, value: number): void {
    const hashkey = this.getHash(key)
    for (const item of this.table[hashkey]) {
      if (item[0] === key) {
        item[1] = value
        return
      }
    }
    this.table[hashkey].push([key, value])
  }

  get(key: number): number {
    const hashkey = this.getHash(key)
    for (const item of this.table[hashkey]) {
      if (item[0] === key) {
        return item[1]
      }
    }
    return -1
  }

  remove(key: number): void {
    const hashkey = this.getHash(key)
    for (const [index, item] of this.table[hashkey].entries()) {
      if (item[0] === key) {
        this.table[hashkey].splice(index, 1)
        return
      }
    }
  }

  // 矩阵的列
  private getHash(key: number) {
    return key % this.bucketsNum
  }
}

export {}
// HashMap 是在 时间和空间 上做权衡的经典例子：
// 如果不考虑空间，我们可以直接设计一个超大的数组，
// 使每个key 都有单独的位置，则不存在冲突；
// 如果不考虑时间，我们可以直接用一个无序的数组保存输入，
// 每次查找的时候遍历一次数组。

// 为了时间和空间上的平衡，HashMap 基于数组实现，
// 并通过 hash 方法求键 key 在数组中的位置，
// 当 hash 后的位置存在冲突的时候，再解决冲突。
