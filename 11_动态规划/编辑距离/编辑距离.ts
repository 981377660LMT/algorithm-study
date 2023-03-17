/**
 * @param {string} word1
 * @param {string} word2
 * @return {number}
 * @description 不优化的diff算法O(n^3)的原因’
 * 传统Diff算法需要找到两个树的最小更新方式，所以需要[两两]对比每个叶子节点是否相同，
 * 对比就需要O(n^2)次了，
 * 再加上更新（移动、创建、删除）时需要遍历一次，所以是O(n^3)。
 * React认为：一个ReactElement的type不同，那么内容基本不会复用，所以直接删除节点，
 * 添加新节点，这是一个非常大的优化，大大减少了对比时间复杂度。
 * @description 你可以对一个单词进行如下三种操作:增删改
 * @summary 编辑距离 O(n*m)
 */
function editDistance<T>(word1: ArrayLike<T>, word2: ArrayLike<T>): number {
  const n1 = word1.length
  const n2 = word2.length
  const dp = new Uint32Array((n1 + 1) * (n2 + 1))

  for (let i = 0; i < n1 + 1; i++) {
    dp[i * (n2 + 1)] = i
  }

  for (let j = 0; j < n2 + 1; j++) {
    dp[j] = j
  }

  for (let i = 1; i < n1 + 1; i++) {
    for (let j = 1; j < n2 + 1; j++) {
      if (word1[i - 1] === word2[j - 1]) {
        dp[i * (n2 + 1) + j] = dp[(i - 1) * (n2 + 1) + j - 1]
      } else {
        dp[i * (n2 + 1) + j] =
          Math.min(
            dp[(i - 1) * (n2 + 1) + j],
            dp[i * (n2 + 1) + j - 1],
            dp[(i - 1) * (n2 + 1) + j - 1]
          ) + 1
      }
    }
  }

  return dp[n1 * (n2 + 1) + n2]
}

if (require.main === module) {
  const s1 = 'horse'.repeat(2000)
  const s2 = 'rosse'.repeat(2000)
  console.time('editDistance')
  console.log(editDistance(s1, s2), s1.length * s2.length)
  console.timeEnd('editDistance')
}

export { editDistance }
