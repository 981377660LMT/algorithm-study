// 2782. 唯一类别的数量
// 有 n 个元素，编号从 0 到 n - 1。每个元素都有一个类别，你的任务是找出唯一类别的数量。
//
// !对每个元素遍历其之前的元素，如果之前的元素有和当前元素相同的类别，那么当前元素不是唯一类别，否则是唯一类别。

declare class CategoryHandler {
  constructor(categories: number[])
  public haveSameCategory(a: number, b: number): boolean
}

function numberOfCategories(n: number, categoryHandler: CategoryHandler): number {
  const visited: number[] = []
  for (let cur = 0; cur < n; cur++) {
    let isUnique = true
    for (let i = 0; i < visited.length; i++) {
      const pre = visited[i]
      if (categoryHandler.haveSameCategory(pre, cur)) {
        isUnique = false
        break
      }
    }

    if (isUnique) {
      visited.push(cur)
    }
  }

  return visited.length
}

export {}
