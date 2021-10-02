// 你可以选择数组的任意下标开始跳跃。请你返回你 最多 可以访问多少个下标。
// 左右 每一步最多为d 且下标要，你从下标 i 跳到下标 j 需要满足：arr[i] > arr[j] 且 arr[i] > arr[k]

// 带记忆化的dfs
// 1 <= arr.length <= 1000
function maxJumps(arr: number[], d: number): number {
  let res = 0
  const memo = new Uint16Array(1001) // 注意初始化为0

  for (let index = 0; index < arr.length; index++) {
    res = Math.max(res, dfs(index))
  }

  return res

  function dfs(cur: number): number {
    if (memo[cur]) return memo[cur]

    let max = 0
    for (let i = cur + 1; i <= cur + d && i < arr.length && arr[i] < arr[cur]; i++) {
      max = Math.max(max, dfs(i))
    }

    for (let i = cur - 1; i >= cur - d && i >= 0 && arr[i] < arr[cur]; i--) {
      max = Math.max(max, dfs(i))
    }

    memo[cur] = max + 1
    return memo[cur]
  }
}

console.log(maxJumps([6, 4, 14, 6, 8, 13, 9, 7, 10, 6, 12], 2))

export {}

// const maxJumps = (arr, d) => {
//   const cache = new Uint16Array(1001);
//   return Math.max(...(arr.map((v, i) => helper(i))));

//   function helper(cur) {
//     if (cache[cur] === 0) {
//       let max = 0;
//       for (let i = cur + 1; i <= cur + d && i < arr.length && arr[i] < arr[cur]; ++i) {
//         max = Math.max(helper(i), max);
//       }
//       for (let i = cur - 1; i >= cur - d && i >= 0 && arr[i] < arr[cur]; --i) {
//         max = Math.max(helper(i), max);
//       }
//       cache[cur] = 1 + max;
//     }
//     return cache[cur];
//   }
// };
