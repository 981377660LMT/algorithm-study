interface NestDict<V> {
  [key: string]: NestDict<V> | V
}

const getNestDictSymbol = (item: NestDict<number | string>) => {
  const res: string[] = []
  const dfs = (cur: NestDict<number | string> | string | number, path: string[]) => {
    if (typeof cur !== 'object') {
      path.push(`#${Object.prototype.toString.call(cur)}:${cur}#`)
      // 为了节省空间，其实这里可以哈希
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

console.log(getNestDictSymbol({ a: { b: 1, c: '2' } }))
