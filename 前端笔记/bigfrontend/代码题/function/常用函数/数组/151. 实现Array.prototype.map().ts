// copied from lib.es5.d.ts
declare interface Array<T> {
  myMap<U>(callbackfn: (value: T, index: number, array: T[]) => U, thisArg?: any): U[]
}

// 怎么获取输入数组的泛型参数
Array.prototype.myMap = function <U>(
  callbackfn: (value: any, index: number, array: any[]) => U,
  thisArg?: any
): U[] {
  const res: U[] = []
  const len = this.length
  for (let i = 0; i < len; i++) {
    // 空值需要忽略
    if (i in this) {
      res[i] = callbackfn.call(thisArg, this[i], i, this)
    }
  }
  return res
}

if (require.main === module) {
  const testArr = Array(2)
  console.log(0 in testArr)
  console.log(1 in testArr)
  console.log(2 in testArr)
}
