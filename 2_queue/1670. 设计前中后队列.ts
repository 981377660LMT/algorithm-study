// 2个双向队列
// 用2个双向队列维护，front维护前一半，back维护后一半
// 类似于数据流的中位数那题

import { ArrayDeque } from './Deque/ArrayDeque'

class FrontMiddleBackQueue {
  private left: ArrayDeque // left不超过right
  private right: ArrayDeque // right不超过left+1

  constructor() {
    this.left = new ArrayDeque(1000)
    this.right = new ArrayDeque(1000)
  }

  pushFront(val: number): void {
    this.left.unshift(val)
    if (this.left.length > this.right.length) {
      // console.log(this.left, 77)
      this.right.unshift(this.left.pop()!)
    }
  }

  pushMiddle(val: number): void {
    if (this.left.length < this.right.length) {
      this.left.push(val)
    } else {
      this.right.unshift(val)
    }
  }

  pushBack(val: number): void {
    if (this.left.length < this.right.length) {
      this.left.push(this.right.shift()!)
    }
    this.right.push(val)
  }

  popFront(): number {
    if (this.left.length < this.right.length) {
      this.left.push(this.right.shift()!)
    }
    return this.left.length ? this.left.shift()! : -1
  }

  popMiddle(): number {
    if (this.left.length === this.right.length) {
      return this.left.length ? this.left.pop()! : -1
    } else {
      return this.right.length ? this.right.shift()! : -1
    }
  }

  popBack(): number {
    if (this.left.length && this.left.length === this.right.length) {
      this.right.unshift(this.left.pop()!)
    }
    return this.right.length ? this.right.pop()! : -1
  }
}

const obj = new FrontMiddleBackQueue()
obj.pushFront(1)
// console.log(obj)
obj.pushBack(2)
obj.pushMiddle(3)
obj.pushMiddle(4)
console.log(obj.popFront())
// 请注意当有 两个 中间位置的时候，选择靠前面的位置进行操作 即选取(index-1)>>1
/**
 * Your FrontMiddleBackQueue object will be instantiated and called as such:
 * var obj = new FrontMiddleBackQueue()
 * obj.pushFront(val)
 * obj.pushMiddle(val)
 * obj.pushBack(val)
 * var param_4 = obj.popFront()
 * var param_5 = obj.popMiddle()
 * var param_6 = obj.popBack()
 */
