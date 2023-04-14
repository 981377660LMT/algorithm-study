/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable @typescript-eslint/no-this-alias */

/**
 * 可持久化栈.
 */
class PersistentStack<T> {
  private readonly _root: SNode<T> | undefined

  constructor(root?: SNode<T>) {
    this._root = root
  }

  push(value: T): PersistentStack<T> {
    return new PersistentStack(new SNode(value, this._root))
  }

  pop(): PersistentStack<T> {
    return new PersistentStack(this._root ? this._root.pre : undefined)
  }

  top(): T | undefined {
    return this._root ? this._root.value : undefined
  }

  empty(): boolean {
    return !this._root
  }

  reverse(): PersistentStack<T> {
    let res = new PersistentStack<T>()
    let x: PersistentStack<T> = this
    while (!x.empty()) {
      res = res.push(x.top()!)
      x = x.pop()
    }
    return res
  }

  toString(): string {
    const sb: T[] = []
    let x: PersistentStack<T> = this
    while (!x.empty()) {
      sb.push(x.top()!)
      x = x.pop()
    }
    sb.reverse()
    return `Stack{${sb.join(', ')}}`
  }
}

class SNode<T> {
  value: T
  pre: SNode<T> | undefined
  constructor(value: T, pre?: SNode<T>) {
    this.value = value
    this.pre = pre
  }
}

export { PersistentStack }

if (require.main === module) {
  let stack = new PersistentStack<number>()
  stack = stack.push(1)
  stack = stack.push(2)
  stack = stack.push(3)

  stack = stack.reverse()
  console.log(stack.toString())

  stack = stack.reverse()
  console.log(stack.toString())
  stack = stack.pop()
  stack = stack.push(4)
  console.log(stack.toString())
  console.log(stack.top())
}
