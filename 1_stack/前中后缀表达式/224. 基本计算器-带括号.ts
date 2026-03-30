const WEIGHT: Record<string, number> = { '(': 0, '+': 1, '-': 1, '*': 2, '/': 2 }
const OPT: Record<string, (a: number, b: number) => number> = {
  '+': (a, b) => a + b,
  '-': (a, b) => a - b,
  '*': (a, b) => a * b,
  '/': (a, b) => Math.trunc(a / b)
}

// 基本计算器-支持+-*/和括号 (双栈+运算符优先级)
// 双栈一算（nums, ops, apply注意b前a后）
// 遇数入数，遇左入左。
// 遇右算到左。
// 遇号先平乱（把栈里面优先级 >= 自己的通通算掉），然后再入栈。
// 最后全算完。
const calculate = function (s: string): number {
  // 在孤立的 +/- 前面补 0，把一元运算变成二元运算
  // 1. 开头是正负号时补 0
  // 2. 左括号后跟正负号时补 0
  // !面试可以暂时不考虑这句
  s = s.replace(/^\s*([+-])/, '0$1').replace(/\(([+-])/g, '(0$1')

  const nums: number[] = []
  const ops: string[] = []
  const apply = () => {
    const b = nums.pop()!
    const a = nums.pop()!
    nums.push(OPT[ops.pop()!](a, b))
  }

  // 数字 | 符号 全局匹配
  for (const part of s.match(/\d+|[+\-*/()]/g) ?? []) {
    if (/\d/.test(part)) {
      nums.push(Number(part))
    } else if (part === '(') {
      ops.push(part)
    } else if (part === ')') {
      while (ops.at(-1) !== '(') apply()
      ops.pop()
    } else {
      // !当前运算符优先级不高于栈顶运算符时，先计算栈顶运算符
      //  eg：1 * 2 + 3
      while (ops.length && ops.at(-1) !== '(' && WEIGHT[ops.at(-1)!] >= WEIGHT[part]) apply()
      ops.push(part)
    }
  }
  while (ops.length) apply()
  return nums[0]
}

console.log(calculate('(1+(4+5+2)-3)+(6+8)'))
// console.log(calculate('2-1+2'))
// console.log(calculate('2147483647'))
// console.log(calculate('-2+ 1'))
// console.log(calculate('1-(-2)'))
console.log(calculate('6-4/2'))

export default 1
