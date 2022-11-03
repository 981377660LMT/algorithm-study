function maxRepeating(sequence: string, word: string): number {
  let left = 1
  let right = Math.floor(sequence.length / word.length)

  // 最右二分
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
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

export {}

// !python:
// return bisect_right(
//   range(1, len(sequence)// len(word) +1 ),
//   False,
//   key=lambda x: word * x not in sequence,
// )
