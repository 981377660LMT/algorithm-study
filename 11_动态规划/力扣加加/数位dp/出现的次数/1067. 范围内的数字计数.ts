import { numberOfKsInRange } from './k出现的次数 copy'

function digitsCount(d: number, low: number, high: number): number {
  return numberOfKsInRange(high, d) - numberOfKsInRange(low - 1, d)
}

console.log(digitsCount(3, 100, 250))
