// 实现一个 sum()，使得如下判断成立。

// const sum1 = sum(1)
// sum1(2) == 3 // true
// sum1(3) == 4 // true
// sum(1)(2)(3) == 6 // true
// sum(5)(-1)(2) == 6 // true

function sum(a: number) {
  const func = (b: number) => sum(a + b)
  //@ts-ignore  两个都可以
  func[Symbol.toPrimitive] = () => a
  func.valueOf = () => a
  return func
}

const sum1 = sum(1)
console.log(+sum1(2))

export {}
