const qs = (arr: number[]): number[] => {
  if (arr.length <= 1) return arr
  const rand = randint(0, arr.length - 1)
  ;[[arr[0], arr[rand]]] = [[arr[rand], arr[0]]]

  const left: number[] = []
  const center: number[] = []
  const right: number[] = []
  const pivot = arr[0]
  arr.forEach(num => {
    if (num < pivot) {
      left.push(num)
    } else if (num === pivot) {
      center.push(pivot)
    } else {
      right.push(pivot)
    }
  })

  return [...qs(left), ...center, ...qs(right)]
}

/**
 * @description 生成[start,end]的随机整数
 */
const randint = (start: number, end: number) => {
  if (start > end) throw new Error('invalid interval')
  const amplitude = end - start
  return Math.floor((amplitude + 1) * Math.random()) + start
}

console.log(qs([9, 4, 10, 3, 1, 1, 0, 10, 8, 3, 9, 9, 4, 10, 10, 9, 9, 9, 1, 0]))
export {}
