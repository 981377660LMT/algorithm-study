function a(msg: string) {
  return msg + 'a'
}
function b(msg: string) {
  return msg + 'b'
}
function c(msg: string) {
  return msg + 'c'
}

const f = compose(a, b, c)
console.log(f('hello'))

type Function = (...args: any[]) => any
function compose(...funcs: Function[]): Function {
  const composeTwo =
    (f1: Function, f2: Function) =>
    (...args: any[]) =>
      f2(f1(...args))
  return funcs.reduce(composeTwo)
}

export {}
