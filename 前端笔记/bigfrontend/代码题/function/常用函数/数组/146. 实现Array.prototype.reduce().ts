declare interface Array<T> {
  myReduce(
    callbackfn: (previousValue: T, currentValue: T, currentIndex: number, array: T[]) => T
  ): T
  myReduce(
    callbackfn: (previousValue: T, currentValue: T, currentIndex: number, array: T[]) => T,
    initialValue: T
  ): T
  myReduce<U>(
    callbackfn: (previousValue: U, currentValue: T, currentIndex: number, array: T[]) => U,
    initialValue: U
  ): U
}

Array.prototype.myReduce = function (...args: any[]) {
  // your code here
  const hasInitialValue = args.length > 1
  if (!hasInitialValue && this.length === 0) {
    throw new Error('Reduce of empty array with no initial value')
  }

  let res = hasInitialValue ? args[1] : this[0]
  for (let i = hasInitialValue ? 0 : 1; i < this.length; i++) {
    res = args[0](res, this[i], i, this)
  }

  return res
}

if (require.main === module) {
  console.log([1, 2, 3].myReduce((pre, cur) => pre + cur, 0))
}

// function (...args: any[]) 用js写实际逻辑
// declare interface Array<T> 用ts欺骗ide
