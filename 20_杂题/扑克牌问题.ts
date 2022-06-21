// 魔术师手中有一堆扑克牌，观众不知道它的顺序，接下来魔术师：

// 从牌顶拿出一张牌， 放到桌子上
// 再从牌顶拿一张牌， 放在手上牌的底部
// 如此往复（不断重复以上两步），直到魔术师手上的牌全部都放到了桌子上。

// 此时，桌子上的牌顺序为： (牌顶) 1,2,3,4,5,6,7,8,9,10,11,12,13 (牌底)。

// 问：原来魔术师手上牌的顺序，用函数实现。
const cal = (result: number[]) => {
  const origin: number[] = []
  while (result.length) {
    if (origin.length) origin.unshift(origin.pop()!)
    origin.unshift(result.pop()!)
  }

  return origin
}

// 正向
const magic = (nums: number[]) => {
  const res: number[] = []
  while (nums.length) {
    res.push(nums.shift()!)
    nums.push(nums.shift()!)
  }
  return res
}
console.log(cal([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13]))

export {}
