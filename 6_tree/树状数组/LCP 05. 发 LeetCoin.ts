class BIT {
  public size: number
  private tree1: number[]
  private tree2: number[]

  constructor(size: number) {
    this.size = size
    // 记录∑di
    this.tree1 = Array(size + 1).fill(0)
    // 记录∑(i-1)*di
    this.tree2 = Array(size + 1).fill(0)
    // a1+a2+...+ar=r*∑di-∑(i-1)*di
  }

  addRange(left: number, right: number, k: number): void {
    this.add(left, k)
    this.add(right + 1, -k)
  }

  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }

  private add(x: number, k: number): void {
    if (x <= 0) throw Error('查询索引应为正整数')
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree1[i] += k
      this.tree1[i] %= MOD
      this.tree2[i] += (x - 1) * k
      this.tree2[i] %= MOD
    }
  }

  // 差分数组(上下车问题):
  // 假设现在有一个原数组a(假设a[0] = 0)，有一个数组d，d[i] = a[i] - a[i-1]，那么
  // a[i] = d[1] + d[2] + .... + d[i]
  // 差分数组 diff[i]，存储的是 res[i] - res[i - 1]；而差分数组 diff[0...i] 的和(树状数组更新/查询就是做求和)，就是 res[i] 的值。
  // 则[1,x]范围的和：a[1] + a[2] + a[3] + ... + a[x] = d[1] + d[1] + d[2] + d[1] + d[2] + d[3] + ... + d[1] + d[2] + d[3] + ... + d[x]
  // =x*(Σd[i]) - ∑(i-1)*d[i] (i从1到x)
  private query(x: number): number {
    let res = 0

    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += ((x * this.tree1[i]) % MOD) - (this.tree2[i] % MOD)
      res = (res + MOD) % MOD
    }

    return res
  }

  private lowbit(x: number) {
    return x & -x
  }
}

// 力扣想进行的操作有以下三种：

// 给团队的一个成员（也可以是负责人）发一定数量的LeetCoin；
// 给团队的一个成员（也可以是负责人），以及他/她管理的所有人（即他/她的下属、他/她下属的下属，……），发一定数量的LeetCoin；
// 查询某一个成员（也可以是负责人），以及他/她管理的所有人被发到的LeetCoin之和。

// https://leetcode-cn.com/problems/coin-bonus/solution/xiao-ai-lao-shi-li-kou-bei-li-jie-zhen-t-rut3/
// https://mp.weixin.qq.com/s?__biz=MzkyMzI3ODgzNQ==&mid=2247483674&idx=1&sn=263595b26950ac60e5bf789839d70c9e&chksm=c1e6cd86f691449062d780b96d9af6d9590a71872ebfee980d5b045b9963714043261027c16b&token=1500097142&lang=zh_CN#rd
// 1. dfs将管理和他管理的人映射到一个区间(这部分很巧妙)[a,b] b表示自身的id
// 2. 树状数组区间update/query
const MOD = 1e9 + 7
function bonus(n: number, leadership: number[][], operations: number[][]): number[] {
  const adjList = Array.from<number, number[]>({ length: n + 1 }, () => [])
  const start = Array<number>(n + 1).fill(0)
  const end = Array<number>(n + 1).fill(0)
  let id = 1

  for (const [u, v] of leadership) {
    adjList[u].push(v)
  }

  dfs(1)

  const res: number[] = []
  const bit = new BIT(n)
  for (const [optType, optId, optValue] of operations) {
    switch (optType) {
      case 1:
        bit.addRange(end[optId], end[optId], optValue)
        break
      case 2:
        bit.addRange(start[optId], end[optId], optValue)
        break
      case 3:
        const queryRes = bit.queryRange(start[optId], end[optId])
        res.push(((queryRes % MOD) + MOD) % MOD)
        break
      default:
        throw new Error('invalid optType')
    }
  }

  return res

  function dfs(cur: number): void {
    start[cur] = id
    for (const next of adjList[cur]) {
      dfs(next)
    }
    end[cur] = id
    id++
  }
}

console.log(
  bonus(
    6,
    [
      [1, 2],
      [1, 6],
      [2, 3],
      [2, 5],
      [1, 4],
    ],
    [
      [1, 1, 500],
      [2, 2, 50],
      [3, 1],
      [2, 6, 15],
      [3, 1],
    ]
  )
)
// 第一次查询时，每个成员得到的LeetCoin的数量分别为（按编号顺序）：500, 50, 50, 0, 50, 0;
// 第二次查询时，每个成员得到的LeetCoin的数量分别为（按编号顺序）：500, 50, 50, 0, 50, 15.
// 输出：[650, 665]
