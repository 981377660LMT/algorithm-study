// This is the interface that allows for creating nested lists.
// You should not implement it, or speculate about its implementation
declare class NestedInteger {
  constructor(value?: number)
  isInteger(): boolean
  getInteger(): number | null
  setInteger(value: number): void
  add(elem: NestedInteger): void
  getList(): NestedInteger[]
}

type NestedArray<T> = Array<NestedArray<T> | T> | T

function deserialize(s: string): NestedInteger {
  const dfs = (input: NestedArray<number> | number): NestedInteger => {
    if (!Array.isArray(input)) {
      return new NestedInteger(input)
    }

    const int = new NestedInteger()
    for (const nextInput of input) {
      int.add(dfs(nextInput))
    }

    return int
  }

  return dfs(JSON.parse(s))
}

console.log(deserialize('[123,[456,[789]]]'))

export {}
