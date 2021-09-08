import { NestedDict } from './10_获取nestdict的唯一识别'

const arr = [0, 'a', 1, 'b', 2, 'c', 3, 'e', 3, 'foo', 2, 'd', 1, 'x', 0, 'ff']

// [0, "a", 1, "b", 2, "c", 3, "e", 2, "d", 1, "x", 0, "ff"] 的一个数组转化为如下的数据。
// {
//   a: {
//     b: {
//       c: {
//         e: null,
//       },
//       d: null,
//     },
//     x: null,
//   },
//   ff: null,
// };

interface StackItem {
  index: number
  value: string
}
//前一题的进阶 需要先将所有的path求出 再用前缀树建树的方法
const entryToTree = (arr: (string | number)[]): NestedDict<null> => {
  const getPathes = (arr: (string | number)[]): string[][] => {
    const pathes: string[][] = []
    const stack: StackItem[] = []
    const entries: StackItem[] = []

    for (let i = 0; i < arr.length; i += 2) {
      entries.push({ index: arr[i] as number, value: arr[i + 1] as string })
    }

    entries.forEach(entry => {
      let shouldPush = true
      while (stack.length && stack[stack.length - 1].index >= entry.index) {
        shouldPush && pathes.push(stack.map(v => v.value))
        shouldPush = false
        stack.pop()
      }
      stack.push(entry)
    })
    stack.length && pathes.push(stack.map(v => v.value))

    return pathes
  }

  const pathes = getPathes(arr)
  console.log(pathes)
  const res = { root: {} } as Record<string, any>
  for (const path of pathes) {
    let root = res.root
    const key = path[path.length - 1]
    const value = null
    for (let i = 0; i < path.length - 1; i++) {
      const char = path[i]
      if (!Object.keys(root).includes(char)) root[char] = {}
      root = root[char]
    }
    root[key] = value
  }

  return res.root
}

console.dir(entryToTree(arr), { depth: null })

export {}
