/**
 * @param {number[]} arr
 * @param {number} start
 * @return {boolean}
 * 当你位于下标 i 处时，你可以跳到 i + arr[i] 或者 i - arr[i]。
 * 请你判断自己能不能跳到值为0 的位置
 * 注意，不管是什么情况下，你都无法跳到数组之外。
 */
const canReach = function (arr: number[], start: number): boolean {
  const dfs = (index: number): boolean => {
    if (index < 0 || index >= arr.length || arr[index] === -1) return false
    if (arr[index] === 0) return true
    const step = arr[index]
    arr[index] = -1
    return dfs(index + step) || dfs(index - step)
  }

  return dfs(start)
}
