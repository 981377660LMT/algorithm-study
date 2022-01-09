function product<T>(arr: T[][]): T[][] {
  const res: T[][] = []
  const target = arr.length
  bt(0, target, [])
  return res

  function bt(i: number, target: number, path: T[]) {
    if (path.length === target) {
      res.push(path.slice())
      return
    }

    for (const choose of arr[i]) {
      path.push(choose)
      bt(i + 1, target, path)
      path.pop()
    }
  }
}

if (require.main === module) {
  console.log(product([['A', 'a'], ['1'], ['B', 'b'], ['2']]))
}

export { product }
