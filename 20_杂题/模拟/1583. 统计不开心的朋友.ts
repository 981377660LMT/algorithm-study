// 排在列表前面的朋友与 i 的亲近程度比排在列表后面的朋友更高
// pairs[i] = [xi, yi] 表示 xi 与 yi 配对，且 yi 与 xi 配对。
function unhappyFriends(n: number, preferences: number[][], pairs: number[][]): number {
  // 每个点x与其他点y的关系排名
  const order = Array.from<unknown, number[]>({ length: n }, () => Array(n).fill(0))
  for (let me = 0; me < n; me++) {
    for (let intimacy = 0; intimacy < n - 1; intimacy++) {
      const friend = preferences[me][intimacy]
      order[me][friend] = intimacy
    }
  }

  // 现在的配对情况
  const match = Array<number>(n).fill(0)
  for (const [cur, next] of pairs) {
    match[cur] = next
    match[next] = cur
  }

  let res = 0
  for (let me = 0; me < n; me++) {
    const myMatch = match[me]
    const ourIntimacy = order[me][myMatch]
    for (let betterIntimacy = 0; betterIntimacy < ourIntimacy; betterIntimacy++) {
      const myBetterMatch = preferences[me][betterIntimacy]
      const matchOfMyBetterMatch = match[myBetterMatch]
      // 她原来更爱我
      if (order[myBetterMatch][me] < order[myBetterMatch][matchOfMyBetterMatch]) {
        res++
        break
      }
    }
  }

  return res
}

console.log(
  unhappyFriends(
    4,
    [
      [1, 2, 3],
      [3, 2, 0],
      [3, 1, 0],
      [1, 2, 0],
    ],
    [
      [0, 1],
      [2, 3],
    ]
  )
)

// 之前男生x和女生u是一对，轰轰烈烈的爱了一场，但是迫于生活和现实，分手走散了。
// 若干年后，x在前女友u的婚礼上，男生x发现，自己迫于无奈找了一个不太爱的人y；而前女友u也找了一个不太爱的人v。
// 男生x酒后痛哭，觉得非常难过，我们最终还是败给了现实。
