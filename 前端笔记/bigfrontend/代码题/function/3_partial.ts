import { Func } from '../typings'

const greet = (greeting: string, name: string) => greeting + ' ' + name + '!'
const greetHello = partial(greet, 'Hello')
greetHello('John') // 'Hello John!'

function partial(fn: Func, ...partials: any[]) {
  return (...args: any[]) => fn(...partials, ...args)
}
// partialRight
function partialRight(fn: Func, ...partials: any[]) {
  return (...args: any[]) => fn(...args, ...partials)
}
