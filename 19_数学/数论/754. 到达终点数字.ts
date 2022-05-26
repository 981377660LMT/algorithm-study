/**
 * @param {number} target
 * @return {number}
 * 转化为对1，2，3，4，5....i,添加正负号，使得其和等于target的最小的数目。
 * S-target=2N,S=(i*(i+1))/2;
   因此就变为找找到最小的i,使得S>=target且S-target为偶数!!!
 */
function reachNumber(target: number): number {
  if (target < 0) target = -target
  if (target === 0) return 0

  let i = 1
  while (true) {
    const sum = (i * (i + 1)) / 2
    if (sum >= target && (sum - target) % 2 === 0) return i
    i++
  }
}

console.log(reachNumber(3))

export {}
