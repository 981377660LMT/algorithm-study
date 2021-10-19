defer(console.log, ['a']), console.log('b') // logs 'b' then 'a'

function defer<TArgs extends any[]>(callback: (...args: TArgs) => void, args: TArgs) {
  setTimeout(callback, 1, ...args)
}
