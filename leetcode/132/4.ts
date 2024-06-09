export {}

const INF = 2e9 // !超过int32使用2e15
// let n = nums.len() as i32;
// let mut memo = vec![0i16; (n * (n + 1) * (k + 1)) as usize];
// (0..memo.len()).for_each(|i| {
//     memo[i] = -1;
// });

//     fn dfs(
//         index: i32,
//         pre: i32,
//         count: i32,
//         n: i32,
//         k: i32,
//         memo: &mut Vec<i16>,
//         nums: &Vec<i32>,
//     ) -> i16 {
//         if count > k {
//             return -1;
//         }
//         if index == n {
//             return 0;
//         }

//         let hash = (index * n * (k + 1) + pre * (k + 1) + count) as usize;
//         if memo[hash] != -1 {
//             return memo[hash];
//         }

//         let mut res = 0;
//         let bad = (pre != 0 && nums[(pre - 1) as usize] != nums[index as usize]) as i32;
//         res = res.max(dfs(index + 1, index + 1, count + bad, n, k, memo, nums) + 1);
//         res = res.max(dfs(index + 1, pre, count, n, k, memo, nums));
//         memo[hash] = res;
//         res
//     }

//     dfs(0, 0, 0, n, k, &mut memo, &nums) as i32
// }

function maximumLength(nums: number[], k: number): number {
  const arr: number[] = []
  const scores: number[] = []
  {
    const n = nums.length
    let ptr = 0
    while (ptr < n) {
      const leader = nums[ptr]
      const start = ptr
      ptr++
      while (ptr < n && nums[ptr] === leader) {
        ptr++
      }
      arr.push(leader)
      scores.push(ptr - start)
    }
  }

  const n = arr.length
  const memo = new Int16Array(n * (n + 1) * (k + 1)).fill(-1)
  const dfs = (index: number, pre: number, count: number): number => {
    if (count > k) {
      return ~n
    }
    if (index === n) {
      return 0
    }
    const hash = index * (n + 1) * (k + 1) + pre * (k + 1) + count
    if (~memo[hash]) {
      return memo[hash]
    }
    let res = 0
    const bad = +(pre && arr[pre - 1] !== arr[index])
    res = Math.max(dfs(index + 1, index + 1, count + bad) + scores[index], res)
    res = Math.max(dfs(index + 1, pre, count), res)
    memo[hash] = res
    return res
  }

  return dfs(0, 0, 0)
}

function bruteFoce(nums: number[], k: number): number {
  const n = nums.length
  const memo = new Int16Array(n * (n + 1) * (k + 1)).fill(-1)
  const dfs = (index: number, pre: number, count: number): number => {
    if (count > k) {
      return ~n
    }
    if (index === n) {
      return 0
    }
    const hash = index * (n + 1) * (k + 1) + pre * (k + 1) + count
    if (~memo[hash]) {
      return memo[hash]
    }
    let res = 0
    const bad = +(pre && nums[pre - 1] !== nums[index])
    res = Math.max(dfs(index + 1, index + 1, count + bad) + 1, res)
    res = Math.max(dfs(index + 1, pre, count), res)
    memo[hash] = res
    return res
  }
  return dfs(0, 0, 0)
}

if (require.main === module) {
  for (let i = 0; i < 10000000; i++) {
    const nums = Array.from({ length: Math.floor(Math.random() * 100) }, () =>
      Math.floor(Math.random() * 1)
    )
    const k = Math.floor(Math.random() * 100)
    if (maximumLength(nums, k) !== bruteFoce(nums, k)) {
      console.log('error', nums, k)
      break
    }
  }
  console.log('pass')
}
