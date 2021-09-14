/**
 * @param {string} senate
 * @return {string}
 * 如果同时存在R和D就继续进行下一轮消灭，轮数直到只剩下R或者D为止！
 * 消灭的策略是，尽量消灭自己后面的对手，因为前面的对手已经使用过权利了，而后序的对手依然可以使用权利消灭自己的同伴！
 * @summary 使用栈模拟
 */
const predictPartyVictory = function (senate: string): string {
  const queue = senate.split('')
  // stack是战场
  const monoStack: string[] = []
  while (queue.length) {
    const head = queue.shift()!
    if (monoStack.length === 0 || monoStack[monoStack.length - 1] === head) monoStack.push(head)
    else queue.push(monoStack.pop()!) // 后面被前面干掉了 前面重新回去
  }
  return monoStack[0] === 'R' ? 'Radiant' : 'Dire'
}

console.log(predictPartyVictory('RDD'))
// "Dire"

console.log(predictPartyVictory('RD'))
// "Radiant"

export default 1
