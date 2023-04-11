import { UnionFindMap } from '../UnionFind'

// 怎么合并？讨论每个数的广义邻居是否在原数组中存在即可

// 增: state | (1 << i)
// 删：a & 1 << i 条件下 state ^ (1 << i) ；增和删可以统一为state ^ (1 << i)
// 替换(删再换广义邻居*号)：a ^ (1 << i) | 1 << 27

// 总结：hard并查集很多都是实际结点通过`虚拟`的广义邻居相连接

// O(26*n)

function groupStrings(words: string[]): number[] {
  const states: number[] = []
  for (const word of words) {
    let state = 0
    for (const char of word) {
      state |= 1 << (char.codePointAt(0)! - 97)
    }
    states.push(state)
  }

  const uf = new UnionFindMap<number>(states)
  const statesSet = new Set<number>(states)

  for (const state of states) {
    for (let i = 0; i < 26; i++) {
      const addOrRemove = state ^ (1 << i)
      if (statesSet.has(addOrRemove)) uf.union(state, addOrRemove)

      // 替换：广义邻居；*用1<<27表示
      if ((state >> i) & 1) {
        const replace = (state ^ (1 << i)) | (1 << 27)
        uf.union(state, replace)
      }
    }
  }

  const groupCounter = new Map<number, number>()
  for (const state of states) {
    const root = uf.find(state)
    groupCounter.set(root, (groupCounter.get(root) ?? 0) + 1)
  }
  return [uf.getPart(), Math.max(...groupCounter.values())]
}

console.log(groupStrings(['a', 'b', 'ab', 'cde']))
console.log(groupStrings(['a', 'ab', 'abc']))
console.log(groupStrings(['b', 'q'])) // [1,2]
console.log(groupStrings(['web', 'a', 'te', 'hsx', 'v', 'k', 'a', 'roh'])) // [5,4]
