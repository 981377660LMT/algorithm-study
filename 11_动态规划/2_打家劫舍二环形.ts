// 第一个房屋和最后一个房屋是紧挨着的
// 分两种情况:偷第一个和不偷第一个
const maxMoney = (money: number[]) => {
  if (money.length < 2) {
    return money[0] || 0
  }

  const dp1 = [money[0]]
  const dp2 = [0, money[1]]

  for (let index = 1; index < money.length - 1; index++) {
    dp1[index] = Math.max(dp1[index - 1], money[index] + dp1[index - 2] || 0)
  }
  for (let index = 2; index < money.length; index++) {
    dp2[index] = Math.max(dp2[index - 1], money[index] + dp2[index - 2])
  }

  return Math.max(dp1.pop()!, dp2.pop()!)
}

console.log(maxMoney([1, 2, 3, 1]))

export {}
