class UnionFind<T> {
  private map: Map<T, T>

  constructor() {
    this.map = new Map()
  }

  union(key1: T, key2: T) {
    const root1 = this.find(key1)
    const root2 = this.find(key2)
    if (root1 !== root2) {
      this.map.set(root1, root2)
    }
    return this
  }

  // 如果key不存在，将返回key自身
  find(key: T) {
    while (this.map.has(key)) {
      key = this.map.get(key)!
    }
    return key
  }
}

// 两个账户都有一些共同的邮箱地址，则两个账户必定属于同一个人。

// 思路:
// 1.使用并查集记录邮箱之间对应的根(邮箱到邮箱的map)
// 2.记录每个邮箱属于谁
// 3.将同一个用户的邮箱维护在一起(每一个用户对应一个根邮箱，根邮箱再对应该用户所有的邮箱)
const accountMerge = (accounts: string[][]) => {
  // 需要合并集合元素(有相同数的集合就和并在一起)，找到rootEmail
  const uf = new UnionFind<string>()
  // 便于从rootEmail找到用户名
  const emailToUserMap = new Map<string, string>()
  // 从rootEmail找到所有邮箱
  const rootEmailToEmail = new Map<string, string[]>()
  const res: string[][] = []

  // 遍历的技巧:加上条件限制，一次做两件事
  for (let i = 0; i < accounts.length; i++) {
    for (let j = 1; j < accounts[i].length; j++) {
      emailToUserMap.set(accounts[i][j], accounts[i][0])
      if (j < accounts[i].length - 1) {
        uf.union(accounts[i][j], accounts[i][j + 1])
      }
    }
  }

  // 构建rootEmailToEmail
  for (const key of emailToUserMap.keys()) {
    const root = uf.find(key)
    if (!rootEmailToEmail.has(root)) rootEmailToEmail.set(root, [])
    rootEmailToEmail.get(root)!.push(key)
  }

  for (const rootEmail of rootEmailToEmail.keys()) {
    res.push([emailToUserMap.get(rootEmail)!, ...rootEmailToEmail.get(rootEmail)!.sort()])
  }

  return res
}

console.log(
  accountMerge([
    ['David', 'David0@m.co', 'David1@m.co'],
    ['David', 'David3@m.co', 'David4@m.co'],
    ['David', 'David4@m.co', 'David5@m.co'],
    ['David', 'David2@m.co', 'David3@m.co'],
    ['David', 'David1@m.co', 'David2@m.co'],
  ])
)

export {}
