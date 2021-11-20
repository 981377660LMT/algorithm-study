function finalPrices(prices: number[]): number[] {
  const monoStack: number[] = []
  const memo: Map<number, number> = new Map()

  for (let i = 0; i < prices.length; i++) {
    // 栈不为空且当前元素大于栈顶元素
    // 说明当前元素是栈顶元素的下一个更大元素
    // while循环表示当前元素是栈中所有已存元素的下一个更大元素
    while (monoStack.length > 0 && prices[i] <= prices[monoStack[monoStack.length - 1]]) {
      memo.set(monoStack.pop()!, prices[i])
    }
    monoStack.push(i)
  }

  return prices.map((price, index) => price - (memo.get(index) ?? 0))
}

console.log(finalPrices([10, 1, 1, 6]))
