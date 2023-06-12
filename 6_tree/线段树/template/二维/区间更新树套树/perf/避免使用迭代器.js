// https://www.zhihu.com/question/58470561/answer/3067263492
// https://runjs.co/s/JdL11j3ZE

function enumerateSubset(nums, callback, copy = false) {
  const n = nums.length
  dfs(0, [])
  function dfs(index, path) {
    if (index === n) {
      callback(copy ? path.slice() : path)
      return
    }
    dfs(index + 1, path)
    path.push(nums[index])
    dfs(index + 1, path)
    path.pop()
  }
}

function* genSubset(nums, copy = false) {
  const n = nums.length
  yield* dfs(0, [])
  function* dfs(index, path) {
    if (index === n) {
      yield copy ? path.slice() : path
      return
    }
    yield* dfs(index + 1, path)
    path.push(nums[index])
    yield* dfs(index + 1, path)
    path.pop()
  }
}

const arr = Array(21)
  .fill(0)
  .map((_, i) => i)

const time1 = performance.now()
enumerateSubset(arr, () => {})
console.log(performance.now() - time1)

const time2 = performance.now()
for (const _ of genSubset(arr)) {
}
console.log(performance.now() - time2)

// !前者比后者快200倍
