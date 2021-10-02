// 一维数组
interface ArrayItem {
  id: number
  name: string
  // pid: number
  children?: ArrayItem[]
}

const arrayToTree = (arr: ArrayItem[]) => {
  const getParentId = (item: ArrayItem): number => {
    const tmp = item.id.toString()
    // Number('')是0
    return Number(tmp.slice(0, -1))
  }

  const adjMap = new Map<number, ArrayItem[]>()
  for (const item of arr) {
    if (item.id === 0) continue
    const pid = getParentId(item)
    if (!adjMap.has(pid)) adjMap.set(pid, [])
    adjMap.get(pid)!.push(item)
  }

  const dfs = (root: ArrayItem) => {
    if (adjMap.has(root.id)) {
      root.children = []
      for (const child of adjMap.get(root.id)!) {
        root.children.push(child)
        dfs(child)
      }
    }
  }

  // 根数组
  const res = [arr[0]]
  res.forEach(dfs)
  return res
}

const jump: ArrayItem[] = [
  { id: 0, name: '总公司' },
  { id: 1, name: '分公司1' },
  { id: 2, name: '分公司2' },
  { id: 11, name: '分公司1-1' },
  { id: 12, name: '分公司1-2' },
  { id: 21, name: '分公司2-1' },
  { id: 111, name: '分公司1-1-1' },
  { id: 112, name: '分公司1-1-2' },
  { id: 121, name: '分公司1-2-1' },
  { id: 122, name: '分公司1-2-2' },
]

console.dir(arrayToTree(jump), { depth: null })
