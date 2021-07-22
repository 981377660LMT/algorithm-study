/**
 * @param {number} n
 * @param {number} k
 * @return {number[][]}
 * @description 组合要去除排列重复的元素?  保持递增顺序即可保证唯一
 * @description 回溯法剪枝
 */
var combine = function (n: number, k: number) {
  const res: number[][] = []

  const bt = (path: number[], index: number) => {
    if (path.length === k) {
      return res.push(path.slice())
    }

    // 这里不一定要遍历到n
    // i到n中至少要有k-path.length个元素 提速10%
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
