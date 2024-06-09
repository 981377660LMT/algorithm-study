mod design_book;
mod example;
mod pattern;
use std::{time::Instant, vec};

mod lc;
mod lib2;
mod std_test;
mod test_traits;

#[derive(Debug)]
struct Person<'a> {
    name: &'a str,
    age: u8,
}

struct Solution {}

impl Solution {
    fn maximum_length(nums: Vec<i32>, k: i32) -> i32 {
        let n = nums.len() as i32;
        let mut memo = vec![0i16; (n * (n + 1) * (k + 1)) as usize];
        (0..memo.len()).for_each(|i| {
            memo[i] = -1;
        });

        // fn dfs(
        //     index: i32,
        //     pre: i32,
        //     count: i32,
        //     n: i32,
        //     k: i32,
        //     memo: &mut Vec<i16>,
        //     nums: &Vec<i32>,
        // ) -> i16 {
        //     if count > k {
        //         return -1;
        //     }
        //     if index == n {
        //         return 0;
        //     }

        //     let hash = (index * n * (k + 1) + pre * (k + 1) + count) as usize;
        //     if memo[hash] != -1 {
        //         return memo[hash];
        //     }

        //     let mut res = 0;
        //     let bad = (pre != 0 && nums[(pre - 1) as usize] != nums[index as usize]) as i32;
        //     res = res.max(dfs(index + 1, index + 1, count + bad, n, k, memo, nums) + 1);
        //     res = res.max(dfs(index + 1, pre, count, n, k, memo, nums));
        //     memo[hash] = res;
        //     res
        // }

        // dfs(0, 0, 0, n, k, &mut memo, &nums) as i32
        let mut dp = vec![-1; (n * (n + 1) * (k + 1)) as usize];
        (0..dp.len()).for_each(|i| {
            dp[i] = -1;
        });

        dp[0] = 0;
        for i in 0..n {
            for j in 0..=i {
                for c in 0..=k {
                    let hash = (i * n * (k + 1) + j * (k + 1) + c) as usize;
                    if dp[hash] == -1 {
                        continue;
                    }

                    let res = dp[hash];
                    let bad = (j != 0 && nums[(j - 1) as usize] != nums[i as usize]) as i32;
                    let hash = (i + 1) * n * (k + 1) + i * (k + 1) + c + bad;
                    dp[hash as usize] = dp[hash as usize].max(res + 1);
                    let hash = (i + 1) * n * (k + 1) + j * (k + 1) + c;
                    dp[hash as usize] = dp[hash as usize].max(res);
                }
            }
        }

        let mut res = 0;
        for i in 0..n {
            for j in 0..=i {
                for c in 0..=k {
                    let hash = (i * n * (k + 1) + j * (k + 1) + c) as usize;
                    res = res.max(dp[hash]);
                }
            }
        }

        res
    }
}
