function largeGroupPositions(s: string): number[][] {
  const res: number[][] = []
  let pre = 0

  for (let cur = 0; cur < s.length; cur++) {
    // 'aaa'的情形：利用越界特性
    if (s[cur] !== s[cur + 1]) {
      if (cur - pre + 1 >= 3) res.push([pre, cur])
      pre = cur + 1
    }
  }

  return res
}

console.log(largeGroupPositions('abbxxxxzzy'))
// 输入：s = "abbxxxxzzy"
// 输出：[[3,6]]
// 解释："xxxx" 是一个起始于 3 且终止于 6 的较大分组。

export {}
