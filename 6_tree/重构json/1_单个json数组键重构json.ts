const foo = ['a', 'b', 'c', 1]

const gen = (arr: any[]): object => {
  const head = arr.shift()!
  return {
    [head]: arr.length <= 1 ? arr[0] : gen(arr),
  }
}

console.dir(gen(foo), { depth: null })

export {}
