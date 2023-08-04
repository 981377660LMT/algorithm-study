// https://leetcode.cn/problems/custom-interval/submissions/
// !不能使用递归，很容易爆栈
// ms = delay + period * count

const alivedInterval = new Set<number>()

/**
 * 注册一个Interval, 返回注册的id.
 */
function customInterval(fn: () => void, delay: number, period: number): number {
  const id = alivedInterval.size
  alivedInterval.add(id)
  _run(id, fn, delay, period)
  return id
}

function customClearInterval(id: number): void {
  alivedInterval.delete(id)
}

function _run(id: number, fn: () => void, delay: number, period: number, count = 0): void {
  setTimeout(() => {
    if (!alivedInterval.has(id)) return
    fn()
    _run(id, fn, delay, period, count + 1)
  }, delay + period * count)
}

export {}
