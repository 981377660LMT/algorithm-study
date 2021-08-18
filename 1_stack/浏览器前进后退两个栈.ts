class Broswer<T> {
  constructor(private stack1: T[] = [], private stack2: T[] = []) {}

  peek() {
    return this.stack1[this.stack1.length - 1]
  }

  open(val: T) {
    this.stack1.push(val)
    return this
  }

  back() {
    this.stack1.length && this.stack2.push(this.stack1.pop()!)
    return this.peek()
  }

  next() {
    this.stack2.length && this.stack1.push(this.stack2.pop()!)
    return this.peek()
  }
}

const broswer = new Broswer()

broswer.open(1).open(2)
console.log(broswer.peek())
console.log(broswer.next())
console.log(broswer.back())
console.log(broswer.back())
