#[allow(dead_code)]
struct Solution;

#[allow(unused_imports)]
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {}
}

/// 双变量问题，枚举右，维护左.
mod two_sum {
    use super::*;
    use std::{collections::HashMap, i32};

    impl Solution {
        #[allow(dead_code)]
        pub fn two_sum(nums: Vec<i32>, target: i32) -> Vec<i32> {
            let mut mp = HashMap::with_capacity(nums.len());
            for (i, v) in nums.iter().enumerate() {
                if let Some(&j) = mp.get(&(target - v)) {
                    return vec![j as i32, i as i32];
                }
                mp.insert(v, i);
            }
            unreachable!();
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn num_identical_pairs(nums: Vec<i32>) -> i32 {
            let mut mp = HashMap::<i32, i32>::with_capacity(nums.len());
            let mut res = 0;
            for v in nums.into_iter() {
                let count = *mp.entry(v).or_default();
                res += count;
                mp.insert(v, count + 1);
            }
            res
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn contains_nearby_duplicate(nums: Vec<i32>, k: i32) -> bool {
            let mut mp = HashMap::with_capacity(nums.len()); // 自动推断类型
            for (i, v) in nums.into_iter().enumerate() {
                // !hashmap insert 如果有值会返回上一次的值, 没有返回none
                // !所以insert在这里是 "有则读取值，没有则插入值" 的作用(setdefault).
                if let Some(j) = mp.insert(v, i) {
                    if (i - j) as i32 <= k {
                        return true;
                    }
                }
            }
            false
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn max_profit(prices: Vec<i32>) -> i32 {
            let mut res = 0;
            let mut pre_min = i32::MAX;
            for v in prices {
                res = res.max(v - pre_min);
                pre_min = pre_min.min(v);
            }
            res
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn max_sum(nums: Vec<i32>) -> i32 {
            let max_digit = |x: i32| -> i32 {
                let mut res = 0;
                let mut cur = x;
                while cur != 0 {
                    res = res.max(cur % 10);
                    cur /= 10;
                }
                res
            };

            let mut mp = HashMap::with_capacity(nums.len());
            let mut res = -1;
            for v in nums {
                let d = max_digit(v);
                if let Some(&pre_max) = mp.get(&d) {
                    res = res.max(pre_max + v);
                    if v > pre_max {
                        mp.insert(d, v);
                    }
                } else {
                    mp.insert(d, v);
                }
            }

            res
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn maximum_sum(nums: Vec<i32>) -> i32 {
            let digit_sum = |x: i32| -> i32 {
                let mut res = 0;
                let mut cur = x;
                while cur != 0 {
                    res += cur % 10;
                    cur /= 10;
                }
                res
            };

            let mut res = -1;
            let mut mp = HashMap::with_capacity(nums.len());
            for v in nums {
                let d = digit_sum(v);
                if let Some(&pre_max) = mp.get(&d) {
                    res = res.max(pre_max + v);
                    if v > pre_max {
                        mp.insert(d, v);
                    }
                } else {
                    mp.insert(d, v);
                }
            }

            res
        }
    }

    // 1679. K 和数对的最大数目
    // 给你一个整数数组 nums 和一个整数 k 。
    // 每一步操作中，你需要从数组中选出和为 k 的两个整数，并将它们移出数组。
    // 返回你可以对数组执行的最大操作数。
    impl Solution {
        #[allow(dead_code)]
        pub fn max_operations(nums: Vec<i32>, k: i32) -> i32 {
            let mut res = 0;
            let mut mp = HashMap::with_capacity(nums.len());
            for v in nums {
                // !需要直接操作map内的值，所以需要get_mut
                if let Some(count) = mp.get_mut(&(k - v)) {
                    if *count > 0 {
                        *count -= 1;
                        res += 1;
                        continue;
                    }
                }
                *mp.entry(v).or_insert(0) += 1;
            }
            res
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn minimum_card_pickup(cards: Vec<i32>) -> i32 {
            let mut res = i32::MAX;
            let mut mp = HashMap::with_capacity(cards.len());
            for (i, v) in cards.into_iter().enumerate() {
                if let Some(j) = mp.insert(v, i as i32) {
                    res = res.min(i as i32 - j + 1);
                }
            }
            if res == i32::MAX {
                -1
            } else {
                res
            }
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn num_pairs_divisible_by60(time: Vec<i32>) -> i32 {
            let mut res = 0;
            let mut mp: HashMap<i32, i32> = HashMap::with_capacity(time.len());
            for v in time {
                let r = v % 60;
                // !or_default() 会返回一个引用，如果没有则插入默认值，行为和python的defaultdict一致
                let count = *mp.entry((60 - r) % 60).or_default();
                res += count;
                *mp.entry(r).or_default() += 1;
            }
            res
        }
    }

    // 四数相加
    impl Solution {
        #[allow(dead_code)]
        pub fn four_sum_count(
            nums1: Vec<i32>,
            nums2: Vec<i32>,
            nums3: Vec<i32>,
            nums4: Vec<i32>,
        ) -> i32 {
            let mut res = 0;
            let mut mp: HashMap<i32, i32> = HashMap::with_capacity(nums1.len() * nums2.len());
            for &x in &nums1 {
                for &y in &nums2 {
                    *mp.entry(x + y).or_insert(0) += 1;
                }
            }
            for &x in &nums3 {
                for &y in &nums4 {
                    if let Some(count) = mp.get(&-(x + y)) {
                        res += count;
                    }
                }
            }
            res
        }
    }
}

// 前缀和
mod presum {
    use super::*;

    #[allow(dead_code)]
    struct NumArray {
        presum: Vec<i32>,
    }
    impl NumArray {
        fn new(nums: Vec<i32>) -> Self {
            let mut presum = vec![0; nums.len() + 1];
            for (i, &v) in nums.iter().enumerate() {
                presum[i + 1] = presum[i] + v;
            }
            Self { presum }
        }

        fn sum_range(&self, left: i32, right: i32) -> i32 {
            self.presum[right as usize + 1] - self.presum[left as usize]
        }
    }

    impl Solution {
        #[allow(dead_code)]
        pub fn vowel_strings(words: Vec<String>, queries: Vec<Vec<i32>>) -> Vec<i32> {
            fn check(c: u8) -> bool {
                c == b'a' || c == b'e' || c == b'i' || c == b'o' || c == b'u'
            }
            let mut presum = vec![0; words.len() + 1];
            for i in 0..words.len() {
                let cur = words[i].as_bytes();
                presum[i + 1] = presum[i]
                    + if check(cur[0]) && check(cur[cur.len() - 1]) {
                        1
                    } else {
                        0
                    };
            }
            queries
                .into_iter()
                .map(|v| presum[v[1] as usize + 1] - presum[v[0] as usize])
                .collect()
        }
    }
}
