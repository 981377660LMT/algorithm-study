export {}

// 在所有长度 ≤ maxLen 的子串中，找字典序最大的那个
function answerString(word: string, numFriends: number): string {
  if (numFriends === 1) {
    return word
  }

  const n = word.length
  const upper = n - numFriends + 1
  let i = 0
  let j = 1
  let k = 0
  while (j + k < n) {
    if (word[i + k] === word[j + k]) {
      k++
    } else if (word[i + k] > word[j + k]) {
      j = j + k + 1
      k = 0
    } else {
      i = Math.max(i + k + 1, j)
      j = i + 1
      k = 0
    }
  }
  const len = Math.min(n - i, upper)
  return word.slice(i, i + len)
}
