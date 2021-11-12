// 如果你猜错了，我会告诉你，你猜测的数字比我选出的数字是大了还是小了。
declare function guess(num: number): number

function guessNumber(n: number): number {
  let l = 1
  let r = n

  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    const state = guess(mid)
    if (state === 0) return mid
    else if (state === -1) r = mid - 1
    else l = mid + 1
  }

  return -1
}
