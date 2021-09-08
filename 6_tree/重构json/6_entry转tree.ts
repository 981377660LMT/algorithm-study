const entry: Dict<string> = {
  'a.b.c.dd': 'abcdd',
  'a.d.xx': 'adxx',
  'a.e': 'ae',
  b: 'kk',
}

// 要求转换成如下对象
// var output = {
//   a: {
//     b: {
//       c: {
//         dd: 'abcdd'
//       }
//     },
//     d: {
//       xx: 'adxx'
//     },
//     e: 'ae'
//   }
// }

type Dict<T> = Record<PropertyKey, T>

// 类似于前缀树
// 前缀树有root 所以这里res也模拟一个root
const entryToTree = (entry: Dict<string>) => {
  const res = { root: {} } as Record<string, any>
  const pathes: string[][] = []

  for (const [k, v] of Object.entries(entry)) {
    pathes.push([...k.split('.'), v])
  }

  for (const path of pathes) {
    let root = res.root
    const key = path[path.length - 2]
    const value = path[path.length - 1]
    for (let i = 0; i < path.length - 2; i++) {
      const char = path[i]
      if (!Object.keys(root).includes(char)) root[char] = {}
      root = root[char]
    }
    root[key] = value
  }

  return res.root
}

console.dir(entryToTree(entry), { depth: null })

export {}
