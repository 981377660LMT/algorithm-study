const factorial = [1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880]

// 康托展开
// 这里主要为了讲解康托展开的思路，实现的算法复杂度为O(n^2)，实际当n很大时，内层循环计算在当前位之后小于当前位的个数可以用
// 树状数组来处理计算，而不用每次都遍历，这样复杂度可以降为O(nlogn)。
// 假设排列数小于10个
function calCantor(arrangement: number, n: number) {
  const str = arrangement.toString()
  let cantorValue = 0

  for (let i = 0; i < n; i++) {
    // 在当前位之后小于其的个数
    let smaller = 0
    for (let j = i + 1; j < n; j++) {
      if (str[j] < str[i]) smaller++
    }

    cantorValue += factorial[n - 1 - i] * smaller
  }

  return cantorValue
}

console.log(calCantor(34152, 5))
console.log(calCantor(1342, 4))

// 逆康托展开
// 第rank大排列，rank从1开始
// 瓶颈在删除的O(n) 使用sortedList可以O(nlogn)
function calArragement(rank: number, n: number) {
  rank--
  const available = Array.from<unknown, number>({ length: n }, (_, i) => i + 1)
  const sb: number[] = []

  for (let i = n; i >= 1; i--) {
    const [div, mod] = [~~(rank / factorial[i - 1]), rank % factorial[i - 1]]
    rank = mod
    sb.push(available[div])
    available.splice(div, 1)
  }

  return sb.join('')
}

console.log(calArragement(62, 5))
console.log(calArragement(1, 5))
