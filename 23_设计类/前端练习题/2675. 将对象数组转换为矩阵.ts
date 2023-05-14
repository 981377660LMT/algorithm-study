// https://leetcode.cn/problems/array-of-objects-to-matrix/
// 1.将每一行的对象转换为Map
// 2.构建矩阵的第一行(所有的keys)
// 3.填充矩阵的剩余行

type Value = string | number | boolean | null

function jsonToMatrix(arr: any[]): Value[][] {
  const rows: Map<string, Value>[] = []

  arr.forEach(item => {
    const rowMap = new Map<string, Value>()
    dfs(item, [], rowMap)
    rows.push(rowMap)
  })

  const allKey = new Set<string>()
  rows.forEach(mp => {
    for (const key of mp.keys()) {
      allKey.add(key)
    }
  })
  const row0 = [...allKey].sort()

  const res: Value[][] = Array(arr.length + 1)
  res[0] = row0
  for (let i = 1; i < res.length; i++) {
    res[i] = Array(allKey.size).fill('')
    for (let j = 0; j < row0.length; j++) {
      const mp = rows[i - 1]
      if (mp.has(row0[j])) {
        res[i][j] = mp.get(row0[j])!
      }
    }
  }
  return res

  function dfs(obj: any, path: string[], rowMap: Map<string, Value>): void {
    // !基本类型
    if (!isObject(obj)) {
      const key = path.join('.')
      rowMap.set(key, obj)
      return
    }

    // !数组
    if (Array.isArray(obj)) {
      obj.forEach((item, index) => {
        path.push(String(index))
        dfs(item, path, rowMap)
        path.pop()
      })
    }

    // !对象
    Object.entries(obj).forEach(([key, val]) => {
      path.push(key)
      dfs(val, path, rowMap)
      path.pop()
    })
  }
}

function isObject(val: unknown): val is object {
  return typeof val === 'object' && val !== null
}

export {}

if (require.main === module) {
  const arr = [{ a: { b: 1, c: 2 } }, { a: { b: 3, d: 4 } }]
  console.log(jsonToMatrix(arr))
}
