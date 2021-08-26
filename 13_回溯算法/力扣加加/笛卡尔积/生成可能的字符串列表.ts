const mapper = {
  '1': ['a', 'b', 'c'],
  '2': ['d', 'e'],
  '3': ['f', 'g', 'h'],
} as { [k: string]: string[] }

const getAll = (num: number) => {
  const str = num.toString()
  const res: string[] = []

  const bt = (str: string, path: string, index: number) => {
    if (path.length === str.length) return res.push(path)
    for (const next of mapper[str[index]]) {
      bt(str, path + next, index + 1)
    }
  }
  bt(str, '', 0)

  return res
}

console.log(getAll(233))
export {}
