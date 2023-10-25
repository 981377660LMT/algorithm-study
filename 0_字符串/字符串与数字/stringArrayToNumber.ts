/**
 * `数字字符串数组`转数字，例如 ['13', '2', '3'] => 1323.
 * 比 Number(arr.join('')) 快.
 * @param numeric 数字字符串.
 */
function stringArrayToNumber(arr: ArrayLike<string>): number {
  let res = 0
  for (let i = 0; i < arr.length; i++) {
    const s = arr[i]
    for (let j = 0; j < s.length; j++) {
      res = res * 10 + +s[j]
    }
  }
  return res
}

export { stringArrayToNumber }

if (require.main === module) {
  // Number(str)

  const arrs: string[][] = Array.from({ length: 1e7 }, () => [])
  for (let i = 0; i < 1e6; i++) {
    for (let j = 1; j < 10; j++) {
      arrs[i].push(String(j))
    }
  }

  console.time('Number')
  arrs.forEach(arr => {
    Number(arr.join('')) // 204.929ms
  })
  console.timeEnd('Number')

  console.time('stringToNumber')
  arrs.forEach(arr => {
    stringArrayToNumber(arr) // 89.234ms
  })
  console.timeEnd('stringToNumber')
}
