const entry: NestDict = {
  a: {
    b: {
      c: {
        dd: 'abcdd',
      },
    },
    d: {
      xx: 'adxx',
    },
    e: 'ae',
  },
}
// 要求转换成如下对象

// var output = {
//   'a.b.c.dd': 'abcdd',
//   'a.d.xx': 'adxx',
//   'a.e': 'ae',
// }

interface NestDict {
  [key: string]: NestDict | string
}

type Dict<T> = Record<PropertyKey, T>

const dfs = (cur: NestDict | string, path: string[], res: Dict<string> = {}) => {
  // 递归终止
  if (typeof cur === 'string') {
    res[path.join('.')] = cur
    return
  }

  Object.keys(cur).forEach(key => {
    path.push(key)
    dfs(cur[key], path, res)
    path.pop()
  })

  return res
}

console.log(dfs(entry, []))

export {}
