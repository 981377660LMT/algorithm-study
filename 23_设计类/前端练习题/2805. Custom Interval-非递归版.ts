// https://leetcode.cn/problems/custom-interval/submissions/

/* eslint-disable no-await-in-loop */
/* eslint-disable prefer-arrow-callback */
// ms = delay + period * count

// !注意 setTimeout 里不能用yield.
// !考虑使用 await sleep(ms)阻塞同步代码(while循环)， 来代替 setTimeout.

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

async function _run(
  id: number,
  fn: () => void,
  delay: number,
  period: number,
  count = 0
): Promise<void> {
  while (true) {
    await sleep(delay + period * count)
    if (!alivedInterval.has(id)) return
    fn()
    count++
  }
}

function sleep(ms = 0): Promise<void> {
  return new Promise(resolve => {
    setTimeout(resolve, ms)
  })
}

export {}
