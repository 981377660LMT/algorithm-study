import { bibfs } from './双向bfs模板'

/**
 *
 * @param start
 * @param end
 * @param bank
 * 请找出能够使起始基因序列变化为目标基因序列所需的最少变化次数
 */
function minMutation(start: string, end: string, bank: string[]): number {
  const bankSet = new Set(bank)
  if (!bankSet.has(end)) return -1
  const gene = ['A', 'T', 'C', 'G']
  return bibfs(start, end, getNextState)

  function getNextState(curState: string): string[] {
    const nextStates: string[] = []

    for (let i = 0; i < curState.length; i++) {
      for (let j = 0; j < gene.length; j++) {
        const nextState = curState.slice(0, i) + gene[j] + curState.slice(i + 1)
        if (nextState === curState || !bankSet.has(nextState)) continue
        nextStates.push(nextState)
      }
    }

    return nextStates
  }
}

console.log(minMutation('AACCGGTT', 'AACCGGTA', ['AACCGGTA']))
// 返回值: 1
