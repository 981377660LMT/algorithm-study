function swapRange(arr: any[], start: number, end: number): void {
  for (let i = start, j = end - 1; i < j; i++, j--) {
    const tmp = arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
  }
}

function rotateLeft(arr: any[], start: number, end: number, step: number): void {
  const n = end - start
  if (n <= 1 || step === 0) return
  if (step >= n) step %= n
  if (step === 0) return
  swapRange(arr, start, start + step)
  swapRange(arr, start + step, end)
  swapRange(arr, start, end)
}

function rotateRight(arr: any[], start: number, end: number, step: number): void {
  const n = end - start
  if (n <= 1 || step === 0) return
  if (step >= n) step %= n
  if (step === 0) return
  swapRange(arr, start, end - step)
  swapRange(arr, end - step, end)
  swapRange(arr, start, end)
}

export { rotateLeft, rotateRight }

if (require.main === module) {
  const arr = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  rotateLeft(arr, 3, 7, 3)
  console.log(arr)
  rotateRight(arr, 3, 7, 3)
  console.log(arr)
  rotateRight(arr, 1, 2, 3)
  console.log(arr)
}
