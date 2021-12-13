function largeGroupPositions(s: string): number[][] {
  const match = [...s.matchAll(/(\w)\1{2,}/g)]
  return match.map(info => [info.index!, info.index! + info[0].length - 1])
}

console.log(largeGroupPositions('abbxxxxzzy'))
// 输入：s = "abbxxxxzzy"
// 输出：[[3,6]]
// 解释："xxxx" 是一个起始于 3 且终止于 6 的较大分组。

export {}
