// 尝试以 O(n log k) 时间复杂度和 O(n) 空间复杂度解决。(堆)

// 一般方法:使用hasmap存储即可(js sort的复杂度是nlog(n))
// 返回的答案应该按单词出现频率由高到低排序。如果不同的单词有相同出现频率，按字母顺序排序。
const topKFrequent = (words: string[], k: number) => {
  const map = new Map<string, number>()
  words.forEach(word => map.set(word, map.get(word)! + 1 || 1))

  return [...map]
    .sort((a, b) => a[0].localeCompare(b[0]))
    .sort((a, b) => b[1] - a[1])
    .slice(0, k)
    .map(item => item[0])
  // .sort((a, b) => a.localeCompare(b))
}

console.log(
  topKFrequent(['the', 'day', 'is', 'sunny', 'the', 'the', 'the', 'sunny', 'is', 'is'], 4)
)

export {}
