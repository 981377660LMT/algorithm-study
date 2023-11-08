function swapRange(arr: Array<unknown>, start: number, end: number): void {
  for (let i = start, j = end - 1; i < j; i++, j--) {
    const tmp = arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
  }
}

export { swapRange, swapRange as reverseRange }

if (require.main === module) {
  const arr = [1, 2, 3, 4, 5]
  swapRange(arr, 0, 4)
  console.log(arr)
}
