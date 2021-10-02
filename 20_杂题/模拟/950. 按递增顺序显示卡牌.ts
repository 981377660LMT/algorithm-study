// 返回能以递增顺序显示卡牌的牌组顺序。
function deckRevealedIncreasing(deck: number[]): number[] {
  const n = deck.length
  deck.sort((a, b) => a - b)
  const res: number[] = []
  while (res.length < n) {
    res.unshift(deck.pop()!)
    res.unshift(res.pop()!)
  }
  // 这一步优点莫名其妙
  res.push(res.shift()!)
  return res
}

console.log(deckRevealedIncreasing([17, 13, 11, 2, 3, 5, 7]))
// 输出：[2,13,3,11,5,17,7]
// 现在要你通过顺序递增的b，反推出你原始的a
// a--->b的操作步骤是这样的：
// 从a数组头拿一个放入b数组尾部，然后把a数组头拿一个放入a数组尾部
// b--->a反推操作就是：
// 从a数组尾部拿一个放入a数组头，然后从b数组尾拿一个放入a数组头
