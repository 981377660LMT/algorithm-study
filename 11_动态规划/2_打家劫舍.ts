const maxMoney = (money: number[]) => {
  if (money.length === 0) {
    return 0
  }

  // 前两项初始值
  const arr = [0, money[0]]
  for (let index = 2; index <= money.length; index++) {
    arr[index] = Math.max(arr[index - 2] + money[index - 1], arr[index - 1])
  }

  return arr
}

console.log(maxMoney([1, 2, 3, 4, 5]))
