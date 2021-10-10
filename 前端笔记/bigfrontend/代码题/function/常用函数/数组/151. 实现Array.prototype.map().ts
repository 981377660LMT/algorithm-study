// copied from lib.es5.d.ts
declare interface Array<T> {
  myMap<U>(callbackfn: (value: T, index: number, array: T[]) => U, thisArg?: any): U[]
}

// 怎么获取输入数组的泛型参数
Array.prototype.myMap = function <U>(
  callbackfn: (value: any, index: number, array: any[]) => U,
  thisArg?: any
): U[] {
  // your code here
  const res: U[] = []
  this.forEach((...args) => {
    const index = args[1]
    res[index] = callbackfn.apply(thisArg, args)
  })
  return res
}
