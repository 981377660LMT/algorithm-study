// 求一个数组相邻差为 1 连续子数组(索引差 1 的同时，值也差 1)的总个数
const countSubArray = (arr: number[]): number => {
  let res = 0
  let pre = 0
  for (let i = 1; i < arr.length; i++) {
    if (arr[i] - arr[i - 1] === 1) pre++
    else pre = 0
    res += pre
  }
  return res
}

export {}
console.log(new Uint8Array(1).fill(255))
console.log(new Uint8Array(1).fill(256))
console.log(new Uint8Array(1).fill(-1)[0])

// push的差异
const a = Array(5)
const b = new Uint32Array(5)
a.push(1)
b.set([1])
console.log(a, b)
type Foo = Exclude<keyof Array<any>, keyof Uint32Array>
