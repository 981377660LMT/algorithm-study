/**
 *
 * @param label
 * 之行排列的计算结果是顺序排列的 对称点
 * 之字变换前后的 label 相加是一个定值
 * 每行最小值是2 ** (level - 1)，最大值是2 ** level - 1，其中 level 是树的深度。
 */
function pathInZigZagTree(label: number): number[] {
  const getLevel = (num: number) => {
    let res = 0
    while (2 ** res - 1 < num) res++
    return res
  }

  const res: number[] = []
  let level = getLevel(label)
  while (level) {
    res.push(label)
    const min = 2 ** (level - 1)
    const max = 2 ** level - 1
    const match = min + max - label
    label = match >> 1
    level--
  }

  return res.reverse()
}

console.log(pathInZigZagTree(14))

// 输出：[1,3,4,14]
