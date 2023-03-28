/* eslint-disable no-param-reassign */

// 遍历分割数/划分数/分割方案

/**
 * 按照字典序降序遍历给定整数 n 的所有可能的正整数分割(和等于 n 的正整数组合).
 * @param n 要进行整数划分的目标值。
 * - n = 50（204226）：12 ms
   - n = 60（966467）：65 ms
   - n = 70（4087968）：230 ms
   - n = 80（15796476）：900 ms
 * @param callbackFn 在每次找到有效的整数划分时被调用，参数是一个表示划分的整数数组。
 * @param lenLimit 限制整数划分的最大长度。-1表示没有限制。
 * @param valLimit 限制整数划分中最大整数的值。-1表示没有限制。
 */
function enumeratePartition(
  n: number,
  callbackFn: (partition: readonly number[]) => void,
  lenLimit = -1,
  valLimit = -1
) {
  const dfs = (partition: number[], sum: number): void => {
    if (sum === n) {
      callbackFn(partition)
      return
    }

    if (lenLimit !== -1 && partition.length === lenLimit) {
      return
    }

    let next = partition.length === 0 ? n : partition[partition.length - 1]
    if (valLimit !== -1 && next > valLimit) {
      next = valLimit
    }

    next = Math.min(next, n - sum)
    partition.push(0)
    for (let x = next; x >= 1; x--) {
      partition[partition.length - 1] = x
      dfs(partition, sum + x)
    }

    partition.pop()
  }

  dfs([], 0)
}

if (require.main === module) {
  console.time('enumeratePartition')
  enumeratePartition(80, partition => {
    // // console.log(partition)
    // console.log(partition)
  })
  console.timeEnd('enumeratePartition')
}

export { enumeratePartition }
