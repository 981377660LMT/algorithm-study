function findString(words: string[], s: string): number {
  let l = 0
  let r = words.length - 1

  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    let mid = (l + r) >> 1
    // 当前位置为""时 移动左或者右
    if (words[mid] === '') {
      if (words[r] === s) return r
      else r--
    } else if (words[mid] === s) {
      return mid
    } else if (words[mid] > s) {
      r = mid - 1
    } else {
      l = mid + 1
    }
  }

  return -1
}

// 哈希:O(n)
// 二分: O(logn) 最坏O(n)
console.log(findString(['at', '', '', '', 'ball', '', '', 'car', '', '', 'dad', '', ''], 'ta'))
