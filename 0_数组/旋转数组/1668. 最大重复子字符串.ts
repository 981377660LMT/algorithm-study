function maxRepeating(sequence: string, word: string): number {
  if (sequence.length < word.length || !sequence.includes(word)) return 0

  let left = 1
  let right = ~~(sequence.length / word.length)

  // 最右二分
  while (left <= right) {
    const mid = (left + right) >> 1
    if (!sequence.includes(word.repeat(mid))) {
      right = mid - 1
    } else {
      left = mid + 1
    }
  }

  return right
}

console.log(maxRepeating('ababc', 'ab'))

// 2
