// 如果是最后一个元素，用 -1 替换。
function replaceElements(arr: number[]): number[] {
  const res = Array<number>(arr.length).fill(0)
  let preMax = -1

  for (let i = arr.length - 1; ~i; i--) {
    res[i] = preMax
    preMax = Math.max(preMax, arr[i])
  }

  return res
}

console.log(replaceElements([17, 18, 5, 4, 6, 1]))
