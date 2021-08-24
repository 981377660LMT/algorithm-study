const foo = ['a', 'b', 'c', 1]

const gen = (arr: (string | number)[], index: number): unknown => {
  if (index === arr.length - 1) return arr[index]
  return { [arr[index]]: gen(arr, index + 1) }
}

console.log(gen(foo, 0))

export {}
