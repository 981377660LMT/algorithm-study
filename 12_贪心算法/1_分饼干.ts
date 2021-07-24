// 局部最优:既能满足孩子，又能消耗最少

/**
 *
 * @param children 孩子的胃口阈值
 * @param biscuits 饼干的胃口阈值
 */
const findContentChildren = (children: number[], biscuits: number[]) => {
  const sortFunc = (a: number, b: number) => a - b
  children.sort(sortFunc)
  biscuits.sort(sortFunc)

  /**
   * @description 满足的孩子数
   */
  let i = 0
  biscuits.forEach(biscuit => {
    if (biscuit >= children[i]) {
      i++
    }
  })

  return i
}

console.log(findContentChildren([1, 2], [1, 2, 3]))
console.log(findContentChildren([4, 1, 2, 3], [1, 2]))

export {}
