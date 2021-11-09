const findMaxConsecutiveOnes = (nums: number[]): number =>
  Math.max(
    ...nums
      .join('')
      .split('0')
      .map(str => str.length)
  )

console.log(findMaxConsecutiveOnes([1, 1, 0, 1, 1, 1]))

export {}
