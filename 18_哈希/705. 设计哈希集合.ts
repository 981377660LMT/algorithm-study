// 10**6数据
// 拉链法

// 2.1 定长拉链数组
// HashSet 设计成一个 M * N 的二维数组
// 时间复杂度：O(1)
// 空间复杂度：O(数据范围)
class MyHashSet1 {
  private bucketsNum: number
  private itemsPerBucket: number
  private table: number[][]

  constructor() {
    this.bucketsNum = 1000
    this.itemsPerBucket = 1001 // 质数
    this.table = Array.from({ length: this.bucketsNum }, () => [])
  }

  add(key: number): void {
    const hashKey = this.getHash(key)
    if (!this.table[hashKey].length) {
      this.table[hashKey] = Array(this.itemsPerBucket).fill(0)
    }
    this.table[hashKey][this.getPosition(key)] = 1
  }

  remove(key: number): void {
    const hashKey = this.getHash(key)
    if (this.table[hashKey].length) {
      this.table[hashKey][this.getPosition(key)] = 0
    }
  }

  contains(key: number): boolean {
    const hashKey = this.getHash(key)
    return this.table[hashKey].length > 0 && this.table[hashKey][this.getPosition(key)] === 1
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

// 2.2 不定长拉链数组
// 拉链会根据分桶中的 key 动态增长，更类似于真正的链表。
// 时间复杂度：O(N/b)，N 是元素个数，b 是桶数。
// 空间复杂度：O(N)
class MyHashSet2 {
  private bucketsNum: number
  private table: number[][]

  constructor() {
    this.bucketsNum = 1009 // 质数个的分桶能让数据更加分散到各个桶中。
    this.table = Array.from({ length: this.bucketsNum }, () => [])
  }

  add(key: number): void {
    const hashKey = this.getHash(key)
    if (this.table[hashKey].includes(key)) {
      return
    }
    this.table[hashKey].push(key)
  }

  remove(key: number): void {
    const hashKey = this.getHash(key)
    const index = this.table[hashKey].indexOf(key)
    if (index === -1) return
    this.table[hashKey].splice(index, 1)
  }

  contains(key: number): boolean {
    const hashKey = this.getHash(key)
    console.log(this.table[hashKey])
    return this.table[hashKey].includes(key)
  }

  // 矩阵的列,第几个桶
  private getHash(key: number) {
    return key % this.bucketsNum
  }
}

export {}
// 定长拉链数组
// 优点：两个维度都可以直接计算出来，查找和删除只用两次访问内存。
// 缺点：需要预知数据范围，用于设计第二个维度的数组大小。
// 不定长拉链数组
// 优点：节省内存，不用预知数据范围；
// 缺点：在链表中查找元素需要遍历
// 发现「不定长拉链数组」法速度最快，
// 大块的内存分配也是需要时间的。因此避免大块的内存分配，也是在节省时间。
