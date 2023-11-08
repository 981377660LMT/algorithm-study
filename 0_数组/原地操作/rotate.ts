import { swapRange } from './swapRange'

function rotateLeft(arr: Array<unknown>, start: number, end: number, step: number): void {
  const n = end - start
  if (step >= n) step %= n
  swapRange(arr, start, start + step)
  swapRange(arr, start + step, end)
  swapRange(arr, start, end)
}

function rotateRight(arr: Array<unknown>, start: number, end: number, step: number): void {
  const n = end - start
  if (step >= n) step %= n
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
}
