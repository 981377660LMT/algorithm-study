/**
 *
 * @param arr
 * @param k
 * @returns 组合可取重复元素
 */
function combinationsWithReplacement<T>(arr: T[], k: number): T[][] {
  const res: T[][] = []

  const bt = (cur: number, path: T[]) => {
    if (path.length === k) return res.push(path.slice())

    for (let i = cur; i < arr.length; i++) {
      path.push(arr[i])
      bt(i, path) // 唯一的区别在此：可取重复的元素
      path.pop()
    }
  }

  bt(0, [])

  return res
}

if (require.main === module) {
  console.log(combinationsWithReplacement([1, 1, 3, 4], 2))
}

export { combinationsWithReplacement }
