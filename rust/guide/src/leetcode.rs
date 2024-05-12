struct Solution {}

const INF: i32 = 1e9 as i32 + 10;

impl Solution {
    pub fn find_permutation(nums: Vec<i32>) -> Vec<i32> {
        if is_sorted(&nums) {
            return nums;
        }

        let n = nums.len() as i32;
        let mut res_cost = INF;
        let mut res = vec![INF];

        for i in 0..n {
            let mut memo = vec![0; (n * (1 << n) * n) as usize];
            for i in 0..memo.len() {
                memo[i] = -1;
            }
            let mut transfer = vec![0; (n * (1 << n) * n) as usize];

            fn hash(n: i32, index: i32, visited: i32, pre: i32) -> i32 {
                index * (1 << n) * n + visited * n + pre
            }

            fn dfs(
                index: i32,
                visited: i32,
                pre: i32,
                memo: &mut Vec<i32>,
                transfer: &mut Vec<i32>,
                nums: &Vec<i32>,
                n: i32,
                first: usize,
            ) -> i32 {
                if index == n {
                    return (pre - nums[first]).abs();
                }
                let hash_ = hash(n, index, visited, pre);
                if memo[hash_ as usize] != -1 {
                    return memo[hash_ as usize];
                }

                let mut res_cost = INF;
                for next in 0..n {
                    if visited & (1 << next) > 0 {
                        continue;
                    }
                    let next_cost = dfs(
                        index + 1,
                        visited | (1 << next),
                        next,
                        memo,
                        transfer,
                        nums,
                        n,
                        first,
                    ) + (pre - nums[next as usize]).abs();
                    if next_cost < res_cost {
                        res_cost = next_cost;
                        let cur_hash = hash(n, index, visited, pre);
                        let next_hash = hash(n, index + 1, visited | (1 << next), next);
                        transfer[cur_hash as usize] = next_hash;
                    }
                }
                memo[hash_ as usize] = res_cost;
                return res_cost;
            }

            let tmp = dfs(1, 1 << i, i, &mut memo, &mut transfer, &nums, n, i as usize);

            if tmp < res_cost {
                res_cost = tmp;
                let mut cur_res = vec![i];
                let mut cur_state = hash(n, 1, 1 << i, i);
                for _ in 1..n {
                    cur_state = transfer[cur_state as usize];
                    cur_res.push(cur_state % n);
                }
                res = cur_res;
            }
        }

        res
    }
}

fn is_sorted<T: PartialOrd>(nums: &Vec<T>) -> bool {
    for i in 1..nums.len() {
        if nums[i] < nums[i - 1] {
            return false;
        }
    }
    true
}
