/**
 * @param {string} num
 * @return {number[]}
 * @description 剪枝情况很多
 */
const splitIntoFibonacci = (num: string): number[] => {
  function* bt(remain: string, path: number[]): Generator<number[]> {
    if (path[path.length - 2]) {
      const cur = path[path.length - 2].toString()
      if (cur.length >= 2 && cur[0] === '0') return path.pop()
    }
    if (path.length > 2) {
      if (path[path.length - 1] !== path[path.length - 2] + path[path.length - 3]) return path.pop()
      if (remain.length === 0) yield path
    }

    for (let i = 0; i < remain.length; i++) {
      const next = remain.slice(0, i + 1)
      path.push(parseInt(next))
      yield* bt(remain.slice(i + 1), path)
    }

    path.pop()
  }

  return bt(num, []).next().value || []
}

console.log(splitIntoFibonacci('1011'))
console.log(splitIntoFibonacci('0123'))
console.log(splitIntoFibonacci('1101111'))
// [123,456,579]
