/**
 * @param {Array<(arg: any) => any>} funcs
 * @return {(arg: any) => any}
 * 为了简单，可以假设传给pipe()的方法都只有一个参数
 */
function pipe(funcs: Array<(arg: any) => any>): (arg: any) => any {
  if (funcs.length === 0) return arg => arg

  const pipeTwo = (func1: Function, func2: Function) => {
    return (...args: any[]) => func2(func1(...args))
  }

  return arg => funcs.reduce(pipeTwo)(arg)
}

function pipe2(this: any, funcs: Array<(arg: any) => any>): (arg: any) => any {
  return arg => funcs.reduce((pre, cur) => cur.call(this, pre), arg)
}

const div = (y: number) => (x: number) => x * y
const plus = (y: number) => (x: number) => x + y
const subtract = (y: number) => (x: number) => x - y
const divide = (y: number) => (x: number) => x / y

console.log(pipe([div(2), plus(3), div(4)])(1)) // 20  // (x * 2 + 3) * 4
console.log(pipe2([div(2), plus(3), div(4)])(1)) // 20  // (x * 2 + 3) * 4
