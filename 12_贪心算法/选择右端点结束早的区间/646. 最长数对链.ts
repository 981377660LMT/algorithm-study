// 有多少不重叠的区间
var findLongestChain = function (pairs: number[][]) {
  // sort by the earliest finish time
  pairs.sort((a, b) => a[1] - b[1])
  let prev = pairs[0],
    chain = 1

  for (let i = 1; i < pairs.length; i++) {
    const [prevS, prevE] = prev
    const [currS, currE] = pairs[i]
    if (prevE < currS) {
      prev = pairs[i]
      chain++
    }
  }
  return chain
}

console.log(
  findLongestChain([
    [1, 2],
    [2, 3],
    [3, 4],
  ])
)
export default 1
