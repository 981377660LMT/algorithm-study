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

function deserialize(s: string): NestedInteger {
  if (!s) return new NestedInteger()
  if (s[0] !== '[') return new NestedInteger(Number(s)) // 整数
  if (s.length <= 2) return new NestedInteger() // 空列表

  const res = new NestedInteger()
  let start = 1
  let level = 0 // level来记录 层次关系，说明s[start:i]里的字符都是某一层的，进行递归就可以了
  for (let i = 1; i < s.length; i++) {
    const isEnd = s[i] === ',' || i === s.length - 1
    // 如果遍历到字符串末尾，没有逗号，所以 需要i == len(s)-1这个判断
    if (level === 0 && isEnd) {
      res.add(deserialize(s.slice(start, i)))
      start = i + 1
    } else if (s[i] === '[') {
      level++
    } else if (s[i] === ']') {
      level--
    }
  }

  return res
}

console.log(deserialize('[123,[456,[789]]]'))

export {}
