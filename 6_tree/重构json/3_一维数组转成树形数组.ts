// 一维数组
interface ArrayItem {
  id: number
  name: string
  pid: number
  children?: ArrayItem[]
}

const arrayToTree = (arr: ArrayItem[]) => {
  const adjMap = new Map<number, ArrayItem[]>()
  arr.forEach(item => {
    if (!adjMap.has(item.pid)) adjMap.set(item.pid, [])
    adjMap.get(item.pid)!.push(item)
  })

  const dfs = (root: ArrayItem) => {
    if (adjMap.has(root.id)) {
      root.children = []
      for (const child of adjMap.get(root.id)!) {
        root.children.push(child)
        dfs(child)
      }
    } else {
      // root.children = undefined
    }
  }

  // 根数组
  const res = arr.filter(item => item.pid === 0)
  res.forEach(dfs)
  return res
}

// const arr: ArrayItem[] = [
//   { id: 0, name: '总公司' },
//   { id: 1, name: '分公司1' },
//   { id: 2, name: '分公司2' },
//   { id: 11, name: '分公司1-1' },
//   { id: 12, name: '分公司1-2' },
//   { id: 21, name: '分公司2-1' },
//   { id: 111, name: '分公司1-1-1' },
//   { id: 112, name: '分公司1-1-2' },
//   { id: 121, name: '分公司1-2-1' },
//   { id: 122, name: '分公司1-2-2' },
// ]
const arr: ArrayItem[] = [
  { id: 1, name: '部门1', pid: 0 },
  { id: 2, name: '部门2', pid: 0 },
  { id: 3, name: '部门3', pid: 0 },
  { id: 4, name: '部门4', pid: 1 },
  { id: 5, name: '部门5', pid: 3 },
  { id: 6, name: '部门6', pid: 2 },
  { id: 7, name: '部门7', pid: 5 },
  { id: 8, name: '部门8', pid: 4 },
  { id: 9, name: '部门9', pid: 1 },
  { id: 10, name: '部门10', pid: 2 },
]

console.dir(arrayToTree(arr), { depth: null })

export {}
