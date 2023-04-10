// 在编号为 i 弹簧处按动弹簧，小球可以弹向 0 到 i-1
// 中任意弹簧或者 i+jump[i] 的弹簧（若 i+jump[i]>=N ，则表示小球弹出了机器）。
// 小球位于编号 0 处的弹簧时不能再向左弹。

import { onlineBfs } from '../../../22_专题/implicit_graph/OnlineBfs-在线bfs'
import { Finder } from '../../../22_专题/implicit_graph/RangeFinder/Finder-fastset'

// !你需要将小球弹出机器。请求出最少需要按动多少次弹簧，可以将小球从编号 0 弹簧弹出整个机器，即向右越过编号 N-1 的弹簧。
function minJump(jump: number[]): number {
  const n = jump.length
  const finder = new Finder(n + 1)
  const dist = onlineBfs(
    n + 1,
    0,
    cur => {
      finder.erase(cur)
    },
    cur => {
      const right = Math.min(cur + jump[cur], n)
      if (finder.has(right)) return right
      return finder.prev(cur)
    }
  )
  return dist[n]
}
