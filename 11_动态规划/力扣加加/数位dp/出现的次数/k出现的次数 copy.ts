/**
 * @param {number} upper n <= 10^9
 * @return {number}  [0,n]有多少个k
 * @description 递归分治
 */
const find = function (upper: number, k: number): number {
  if (upper <= k - 1) return 0
  if (upper < 10) return 1

  // 举例:upper=124,k=3
  const [div, mod] = [~~(upper / 10), upper % 10]
  // 1. 先算k在个位出现的次数：
  // 看前面等不等于div
  // 不等于div:个位固定，前面变化有div种(0-div-1)
  // 等于div:find(mod, k)
  let res = div + find(mod, k)
  // 2. k在十位数及以上部分x出现的次数：
  // 看前面等不等于div
  // 不等于div:0-11(0-119) 此时个位数可取0-9 10位 除去前导0 所以是(find(div - 1, k) - find(0, k)) * 10
  // 等于div:12(120-124) 此时个位数可取(mod + 1)位 所以是(find(div, k) - find(div - 1, k)) * (mod + 1)
  res += (find(div - 1, k) - find(0, k)) * 10 + (find(div, k) - find(div - 1, k)) * (mod + 1)
  return res
}

if (require.main === module) {
  console.log(find(251, 2))
  console.log(find(10, 0))
}

export { find as numberOfKsInRange }
