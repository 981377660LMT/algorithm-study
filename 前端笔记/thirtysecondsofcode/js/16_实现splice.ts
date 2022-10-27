const names = ['alpha', 'bravo', 'charlie']
const namesAndDelta = splice(names, 1, 0, 'delta')
// [ 'alpha', 'delta', 'bravo', 'charlie' ]
const namesNoBravo = splice(names, 1, 1) // [ 'alpha', 'charlie' ]
console.log(names) // ['alpha', 'bravo', 'charlie']

// function shank<T>(arr: T[], start: number, deleteCount?: number): T[]
// splice 原地修改
function splice<T>(arr: T[], start: number, deleteCount: number, ...items: T[]): T[] {
  const preLen = arr.length + items.length
  // start右侧的arr 删除之后应添加的数组
  const rightArr = items.concat(arr.slice(start + deleteCount))

  let i = start
  while (rightArr.length) {
    arr[i] = rightArr.shift()!
    i++
  }

  arr.length = preLen - deleteCount

  return arr
}

export {}
