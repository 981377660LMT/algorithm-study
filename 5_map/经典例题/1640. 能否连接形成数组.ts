// 请你以 任意顺序 连接 pieces 中的数组以形成 arr 。但是，不允许 对每个数组 pieces[i] 中的整数重新排序。
// 组间可以换位()，组内不能换位
// 数组中的每个整数 互不相同

// 总结：
// 组间可以换位=>用哈希表查询
// 组内不能换位=>当成一个整体，用首部元素作为key
function canFormArray(arr: number[], pieces: number[][]): boolean {
  const trunkRecord = new Map<number, number[]>()
  const res: number[][] = []

  for (const piece of pieces) {
    trunkRecord.set(piece[0], piece)
  }

  for (const num of arr) {
    if (trunkRecord.has(num)) {
      res.push(trunkRecord.get(num)!)
    }
  }

  console.log(trunkRecord, res)
  return arr.join('#') === res.flat().join('#')
}

console.log(canFormArray([91, 4, 64, 78], [[78], [4, 64], [91]]))
// 输出：true
// 解释：即便数字相符，也不能重新排列 pieces[0]

// 输入：arr = [49,18,16], pieces = [[16,18,49]]
// 输出：false
// 解释：即便数字相符，也不能重新排列 pieces[0]
