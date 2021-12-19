function kIncreasing(arr: number[], k: number): number {
  let res = 0
  for (let start = 0; start < k; start++) {
    const slice: number[] = []
    for (let i = start; i < arr.length; i += k) slice.push(arr[i])
    res += helper(slice)
  }

  return res
}

function helper(arr: number[]): number {
  const LIS = [arr[0]]
  for (let i = 1; i < arr.length; i++) {
    const num = arr[i]
    if (num >= LIS[LIS.length - 1]) {
      LIS.push(num)
    } else {
      const index = bisectRight(LIS, num)
      LIS[index] = num
    }
  }

  return arr.length - LIS.length
}

function bisectRight(arr: number[], target: number): number {
  if (arr.length === 0) return 0

  let l = 0
  let r = arr.length - 1

  while (l <= r) {
    const mid = (l + r) >> 1
    const midElement = arr[mid]
    if (midElement <= target) {
      l = mid + 1
    } else {
      r = mid - 1
    }
  }

  return l
}
