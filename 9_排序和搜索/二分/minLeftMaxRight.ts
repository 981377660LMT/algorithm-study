// minLeftMaxRight

/**
 * 返回最大的 `right` 使得 `[left,right)` 内的值满足 {@link check}.
 *
 * @param upper `right` 的上界, 即 `right<=upper`.
 */
function maxRight(left: number, check: (right: number) => boolean, upper: number): number {
  let ok = left
  let ng = upper + 1
  while (ok + 1 < ng) {
    const mid = Math.floor((ok + ng) / 2)
    if (check(mid)) {
      ok = mid
    } else {
      ng = mid
    }
  }
  return ok
}

/**
 * 返回最小的 `left` 使得 `[left,right)` 内的值满足 {@link check}.
 *
 * @param lower `left` 的下界, 即 `left>=lower`.
 */
function minLeft(right: number, check: (left: number) => boolean, lower: number): number {
  let ok = right
  let ng = lower - 1
  while (ng + 1 < ok) {
    const mid = Math.floor((ok + ng) / 2)
    if (check(mid)) {
      ok = mid
    } else {
      ng = mid
    }
  }
  return ok
}

export { minLeft, maxRight }
