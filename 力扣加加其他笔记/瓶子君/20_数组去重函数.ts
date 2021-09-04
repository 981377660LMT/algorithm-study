// 不能使用JSON.stringify 因为key的位置可能不同

interface NestDict<V> {
  [key: string]: NestDict<V> | V
}

type ArrayItem = number | string | NestDict<number | string>

const unique = (arr: ArrayItem[]) => {
  const getNestDictSymbol = (item: NestDict<number | string>) => {
    const res: string[] = []

    const dfs = (cur: NestDict<number | string> | string | number, path: string[]) => {
      if (typeof cur !== 'object') {
        path.push(`#${Object.prototype.toString.call(cur)}:${cur}#`)
        res.push(path.join(''))
        return path.pop()
      }

      Object.keys(cur)
        .sort()
        .forEach(key => {
          path.push(key)
          dfs(cur[key], path)
          path.pop()
        })
    }

    dfs(item, [])
    return res.join('')
  }

  const res: ArrayItem[] = []
  const visited = new Set<string | number>()
  for (const item of arr) {
    if (typeof item !== 'object' || item === null) {
      if (visited.has(item)) continue
      res.push(item)
      visited.add(item)
    } else {
      const symbol = getNestDictSymbol(item)
      if (visited.has(symbol)) continue
      res.push(item)
      visited.add(symbol)
    }
  }

  return res
}

console.dir(
  unique([123, { a: 1 }, { a: '1' }, 'ok', { a: { b: 1, c: 2 } }, { a: { c: 2, b: 1 } }]),
  {
    depth: null,
  }
)

export {}
