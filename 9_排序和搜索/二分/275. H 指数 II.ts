// 计算并返回该研究者的 h 指数。citations 已经按照 升序排列
// 总共有 h 篇论文分别被引用了至少 h 次。且其余的 n - h 篇论文每篇被引用次数 不超过 h 次。
// h 指数 是其中最大的那个。
function hIndex(citations: number[]): number {
  let [l, r] = [0, citations.length - 1]

  // while (l <= r) {
  //   const mid = (l + r) >> 1
  //   if (citations[mid] > mid) l = mid + 1
  //   else r = mid - 1
  // }

  return l
}

console.log(hIndex([0, 1, 3, 5, 6]))
