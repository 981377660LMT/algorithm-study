/**
 * @param {number[][]} board
 * @return {number}
 * @description 一次移动定义为选择 0 与一个相邻的数字（上下左右）进行交换.
 * @description 最终当板 board 的结果是 [[1,2,3],[4,5,0]] 谜板被解开。
 * @description 实际上六位数比字符串快
 * @description 在这里记录0所在的index 以及index对应下一步的状态
 * @summary 时间复杂度O(n!)
 */
const slidingPuzzle = (board: number[][]): number => {
  let res = -1
  const initState = board.map(row => row.join('')).reduce((pre, cur) => pre + cur, '')
  const queue: [string, number, number][] = [[initState, initState.indexOf('0'), 0]]
  const visited = new Set([initState])
  const nextIndexMapping: { [k: number]: number[] } = {
    0: [1, 3],
    1: [0, 2, 4],
    2: [1, 5],
    3: [0, 4],
    4: [1, 3, 5],
    5: [2, 4],
  }
  const getNextState = (curState: string, curIndex: number, nextIndex: number): string => {
    const tmp = curState.split('')
    ;[tmp[curIndex], tmp[nextIndex]] = [tmp[nextIndex], tmp[curIndex]]
    return tmp.join('')
  }

  const bfs = () => {
    while (queue.length) {
      const [curState, curIndex, steps] = queue.shift()!
      if (curState === '123450') return (res = steps)
      for (const nextIndex of nextIndexMapping[curIndex]) {
        const nextState = getNextState(curState, curIndex, nextIndex)
        if (!visited.has(nextState)) {
          visited.add(nextState)
          queue.push([nextState, nextIndex, steps + 1])
        }
      }
    }
  }
  bfs()

  return res
}

console.log(
  slidingPuzzle([
    [4, 1, 2],
    [5, 0, 3],
  ])
)

export {}
