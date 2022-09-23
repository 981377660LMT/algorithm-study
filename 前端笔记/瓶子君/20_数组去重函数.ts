// 不能使用JSON.stringify 因为key的位置可能不同
// 深度数组去重

interface NestDict<V> {
  [key: string]: NestDict<V> | V
}

type ArrayItem = number | string | NestDict<number | string>

function genHash(item: NestDict<number | string>) {
  const res: string[] = []

  dfs(item, [])
  return res.join('')

  function dfs(cur: NestDict<number | string> | string | number, path: string[]): void {
    if (typeof cur !== 'object') {
      path.push(`#${Object.prototype.toString.call(cur)}:${cur}#`)
      res.push(path.join(''))
      path.pop()
      return
    }

    Object.keys(cur)
      .sort()
      .forEach(key => {
        path.push(key)
        dfs(cur[key], path)
        path.pop()
      })
  }
}

function unique(arr: ArrayItem[]) {
  const res: ArrayItem[] = []
  const visited = new Set<string | number>()
  for (const item of arr) {
    if (typeof item !== 'object' || item == null) {
      if (visited.has(item)) {
        continue
      }
      res.push(item)
      visited.add(item)
    } else {
      const hash = genHash(item)
      if (visited.has(hash)) {
        continue
      }
      res.push(item)
      visited.add(hash)
    }
  }

  return res
}

console.dir(
  unique([123, { a: 1 }, { a: '1' }, 'ok', { a: { b: 1, c: 2 } }, { a: { c: 2, b: 1 } }]),
  {
    depth: null
  }
)

export {}
