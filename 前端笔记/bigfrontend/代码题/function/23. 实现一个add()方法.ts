// 题目描述:实现一个 add 方法 使计算结果能够满足如下预期：
// add(1)(2)(3)()=6 add(1,2,3)(4)()=10

function sum(...curArgs: number[]) {
  const res = (...nextArgs: number[]) =>
    sum([...curArgs, ...nextArgs].reduce((pre, cur) => pre + cur, 0))
  res.valueOf = () => curArgs.reduce((pre, cur) => pre + cur, 0)
  return res
}

const sum1 = sum(1)
console.log(+sum1(1)(2)()(4, 5, 6))

export {}
