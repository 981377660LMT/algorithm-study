function minimumDifference(nums: number[]): number {
  const sum = nums.reduce((pre, cur) => pre + cur, 0)
  let res = Infinity

  nums.sort((a, b) => a - b)

  const bt = (index: number, pathSum: number, visited: number, len: number): void => {
    if (len === nums.length / 2) {
      res = Math.min(res, Math.abs(sum - 2 * pathSum))
      return
    }

    for (let i = index; i < nums.length; i++) {
      if (i !== index && nums[i] === nums[i - 1]) continue
      if (visited & (1 << i)) continue
      visited |= 1 << i
      bt(i + 1, pathSum + nums[i], visited, len + 1)
    }
  }
  bt(0, 0, 0, 0)

  return res
}
console.log(minimumDifference([2, -1, 0, 4, -2, -9]))
console.log(minimumDifference([3, 9, 7, 3]))
console.log(minimumDifference([-36, 36]))
// console.log(minimumDifference([2, -1, 0, 4, -2, -9]))

// interface Stock {
//   price: number
//   timeStamp: number
// }

// class StockPrice {
//   private queue: number[]
//   private maxQueue: number[] // 队头最大
//   private minQueue: number[]

//   constructor() {
//     this.queue = []
//     this.maxQueue = []
//     this.minQueue = []
//   }

//   update(timestamp: number, price: number): void {}

//   current(): number {
//     return this.queue[this.queue.length - 1]
//   }

//   maximum(): number {
//     return this.maxQueue[0]
//   }

//   minimum(): number {
//     return this.minQueue[0]
//   }
// }

// class MaxQueue {
//   private queue: number[]
//   private maxQueue: number[] // 队头最大
//   constructor() {
//     this.queue = []
//     this.maxQueue = []
//   }

//   max_value(): number {
//     return this.maxQueue[0] || -1
//   }

//   push_back(value: number): void {
//     this.queue.push(value)
//     while (this.maxQueue.length && this.maxQueue[this.maxQueue.length - 1] < value) {
//       this.maxQueue.pop()
//     }
//     this.maxQueue.push(value)
//   }

//   pop_front(): number {
//     if (!this.queue.length) return -1
//     const value = this.queue.shift()!
//     value === this.maxQueue[0] && this.maxQueue.shift() // 最大值出队
//     return value
//   }
// }

export {}
