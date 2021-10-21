/**
 * @param {number} n
 * @param {number} k
 * @return {number[][]}
 * @description 组合要去除排列重复的元素?  保持递增顺序即可保证唯一
 * @description 回溯法剪枝
 * 返回范围 [1, n] 中所有可能的 k 个数的组合。
 */
var combine = function (n: number, k: number) {
  const res: number[][] = []

  const bt = (path: number[], index: number) => {
    if (path.length === k) {
      return res.push(path.slice())
    }

    // 这里不一定要遍历到n
    // i 大于n - (k - path.length) + 1则之后path装不满k个数了

    for (let i = index + 1; i <= n - (k - path.length) + 1; i++) {
      path.push(i)
      bt(path, i)
      path.pop()
    }
  }
  bt([], 0)

  return res
}

console.log(combine(4, 2))

export {}
