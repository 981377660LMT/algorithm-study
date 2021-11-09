/**
 *
 * @param label
 * 之行排列的计算结果是顺序排列的 对称点
 * 之字变换前后的 label 相加是一个定值
 * 每行最小值是2 ** (level - 1)，最大值是2 ** level - 1，其中 level 是树的深度。
 */
function pathInZigZagTree(label: number): number[] {
  const res: number[] = []
  let level = getLevel(label)

  while (level) {
    res.push(label)
    const min = 2 ** (level - 1)
    const max = 2 ** level - 1
    const match = min + max - label // 对称
    label = match >> 1
    level--
  }

  return res.reverse()

  function getLevel(num: number) {
    return ~~Math.log2(num) + 1
  }
}

console.log(pathInZigZagTree(14))

// 输出：[1,3,4,14]
