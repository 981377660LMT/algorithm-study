/**
 * @param {number[][]} docs
 * @return {string[]}
 * 两个(具有不同单词的)文档的交集(intersection)中元素的个数除以并集(union)中元素的个数，就是这两个文档的相似度
 */
const computeSimilarities = function (docs: number[][]): string[] {
  // 反向索引 key为关键词  value 为文档的集合
  const map = new Map<number, number[]>()
  for (let i = 0; i < docs.length; i++) {
    for (let j = 0; j < docs[i].length; j++) {
      !map.has(docs[i][j]) && map.set(docs[i][j], [])
      map.get(docs[i][j])!.push(i)
    }
  }
  const valid = [...map.entries()].filter(entry => entry[1].length >= 2)
  const pairMap = new Map<string, number>()
  console.log(valid)
  for (const entry of valid) {
    for (let i = 0; i < entry[1].length - 1; i++) {
      for (let j = i + 1; j < entry[1].length; j++) {
        const key = `${entry[1][i]}#${entry[1][j]}`
        pairMap.set(key, (pairMap.get(key) || 0) + 1)
      }
    }
  }

  console.log(pairMap)
  return [...pairMap.entries()].map(entry => {
    const [key1, key2] = entry[0].split('#')
    const count = entry[1]
    return `${key1},${key2}: ${(count / (docs[+key1].length + docs[+key2].length - count)).toFixed(
      4
    )}`
  })
}

console.log(
  computeSimilarities([
    [14, 15, 100, 9, 3],
    [32, 1, 9, 3, 5],
    [15, 29, 2, 6, 8, 7],
    [7, 10],
  ])
)
// 输出: ['0,1: 0.2500', '0,2: 0.1000', '2,3: 0.1429']
