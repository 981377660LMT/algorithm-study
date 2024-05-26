use std::time::Instant;

mod lc;
mod lib2;
mod std_test;
mod test_traits;

#[derive(Debug)]
struct Person<'a> {
    name: &'a str,
    age: u8,
}

fn main() {
    let name = "Peter";
    let age = 27;
    let peter = Person { name, age };
    println!("{:?}", peter);

    // 专门用于调试的断言
    debug_assert_eq!(peter.age, 27);

    let now = Instant::now();
    Solution::find_permutation(vec![0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13]);
    println!("{:?}", now.elapsed());
    let now = Instant::now();
    let mut sum: i64 = 0;
    for i in 0..1e9 as i32 {
        sum += i as i64;
    }
    println!("{}", sum);
    println!("{:?}", now.elapsed());

    debug_assert_eq!(1, 2);

    let cond = if 1 > 0 { 1 } else { 0 };
} // Person { name: "Peter", age: 27 }

// package main

// import "sort"

// const INF int32 = 1e9 + 10

// func findPermutation(nums []int) []int {
// 	if sort.IntsAreSorted(nums) {
// 		return append([]int(nil), nums...)
// 	}
// 	n := int32(len(nums))
// 	newNums := make([]int32, n)
// 	for i := range nums {
// 		newNums[i] = int32(nums[i])
// 	}

// 	resCost, res := INF, []int32{INF}
// 	for i := int32(0); i < n; i++ {
// 		memo := make([]int32, n*(1<<n)*n)
// 		for i := range memo {
// 			memo[i] = -1
// 		}
// 		next_ := make([]int32, n*(1<<n)*n)
// 		hash := func(index, visited, pre int32) int32 {
// 			return index*(1<<n)*n + visited*n + pre
// 		}

// 		first := i
// 		var dfs func(index, visited, pre int32) int32
// 		dfs = func(index, visited, pre int32) int32 {
// 			if index == n {
// 				return abs(pre - newNums[first])
// 			}
// 			hash_ := hash(index, visited, pre)
// 			if memo[hash_] != -1 {
// 				return memo[hash_]
// 			}

// 			resCost := INF
// 			for next := int32(0); next < n; next++ {
// 				if visited&(1<<next) > 0 {
// 					continue
// 				}
// 				nextCost := dfs(index+1, visited|(1<<next), next) + abs(pre-newNums[next])
// 				if nextCost < resCost {
// 					resCost = nextCost
// 					curHash := hash(index, visited, pre)
// 					nextHash := hash(index+1, visited|(1<<next), next)
// 					next_[curHash] = nextHash
// 				}
// 			}
// 			memo[hash_] = resCost
// 			return resCost
// 		}
// 		tmp := dfs(1, 1<<i, i)

// 		if tmp < resCost {
// 			resCost = tmp
// 			curRes := []int32{i}
// 			curState := hash(1, 1<<i, i)
// 			for i := int32(1); i < n; i++ {
// 				curState = next_[curState]
// 				curRes = append(curRes, curState%n)
// 			}
// 			res = curRes
// 		}
// 	}

// 	newRes := make([]int, n)
// 	for i := range res {
// 		newRes[i] = int(res[i])
// 	}
// 	return newRes
// }

// func abs(a int32) int32 {
// 	if a < 0 {
// 		return -a
// 	}
// 	return a
// }

const INF: i32 = 1e9 as i32 + 10;

struct Solution;

impl Solution {
    pub fn find_permutation(nums: Vec<i32>) -> Vec<i32> {
        let n = nums.len() as i32;
        let mut res_cost = INF;
        let mut res = vec![INF];
        let mut memo = vec![-1; (n * (1 << n) * n) as usize];
        let mut transfer = vec![0; (n * (1 << n) * n) as usize];
        for i in 0..n {
            for i in 0..memo.len() {
                memo[i] = -1;
            }
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
