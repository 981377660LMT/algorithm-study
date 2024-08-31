// 模式匹配中的变量绑定  x @ (0 | 1 | 8) => res = res * 10 + x,
// fuse : fuse() -> Fuse<Self>，创建一个迭代器，它会在迭代器被耗尽后返回 None (保险丝: fuse)
mod dfa;
mod eval_lisp;
mod my_atoi;

struct Solution;

impl Solution {
    pub fn buddy_strings(s: String, goal: String) -> bool {
        if s.len() != goal.len() {
            return false;
        }
        let diff = s
            .bytes()
            .zip(goal.bytes())
            .filter(|u| u.0 != u.1)
            .collect::<Vec<_>>();

        match diff.len() {
            2 => diff[0].0 == diff[1].1 && diff[0].1 == diff[1].0,
            0 => s.bytes().collect::<std::collections::HashSet<_>>().len() < s.len(),
            _ => false,
        }
    }

    pub fn find_missing_ranges(mut nums: Vec<i32>, lower: i32, upper: i32) -> Vec<Vec<i32>> {
        fn find(a: i32, b: i32) -> Option<Vec<i32>> {
            match b - a {
                0 | 1 => None,
                _ => Some(vec![a + 1, b - 1]),
            }
        }

        let mut res = vec![];
        nums.insert(0, lower - 1);
        nums.push(upper + 1);
        for window in nums.windows(2) {
            if let Some(s) = find(window[0], window[1]) {
                res.push(s);
            }
        }

        res
    }

    pub fn confusing_number(n: i32) -> bool {
        let raw_n = n;
        let mut n = n;
        let mut res = 0;
        while n > 0 {
            match n % 10 {
                6 => res = res * 10 + 9,
                9 => res = res * 10 + 6,
                x @ (0 | 1 | 8) => res = res * 10 + x,
                _ => return false,
            }
            n /= 10;
        }
        res != raw_n
    }

    // 485. 最大连续 1 的个数
    // https://leetcode.cn/problems/max-consecutive-ones/description/
    pub fn find_max_consecutive_ones(nums: Vec<i32>) -> i32 {
        nums.split(|&x| x == 0)
            .map(|v| v.iter().count())
            .max()
            .unwrap_or(0) as i32
    }

    pub fn simplify_path(path: String) -> String {
        let paths = path
            .split('/')
            .filter(|&x| x != "." && !x.is_empty())
            .fold(vec![], |mut acc, x| {
                match x {
                    ".." => {
                        acc.pop();
                    }
                    _ => acc.push(x),
                }
                acc
            })
            .join("/");

        "/".to_string() + &paths
    }

    pub fn compare_version(version1: String, version2: String) -> i32 {
        let mut s1 = version1.split(".").fuse().map(|s| s.parse().unwrap());
        let mut s2 = version2.split(".").fuse().map(|s| s.parse().unwrap());
        loop {
            match (s1.next(), s2.next()) {
                (None, None) => break,
                (Some(0), None) | (None, Some(0)) => continue,
                (Some(_n1), None) => return 1,
                (None, Some(_n2)) => return -1,
                (Some(n1), Some(n2)) if n1 > n2 => return 1,
                (Some(n1), Some(n2)) if n1 < n2 => return -1,
                (Some(_n1), Some(_n2)) /*if n1 == n2*/ => continue,
            }
        }
        0
    }

    pub fn longest_valid_parentheses(s: String) -> i32 {
        let mut stack = vec![-1];
        let mut res = 0;
        s.into_bytes().into_iter().enumerate().for_each(|(i, c)| {
            if c == b'(' {
                stack.push(i as i32);
            } else {
                stack.pop();
                if stack.is_empty() {
                    stack.push(i as i32);
                } else {
                    res = res.max(i as i32 - stack.last().unwrap());
                }
            }
        });

        res
    }

    // 78. 子集
    pub fn subsets(nums: Vec<i32>) -> Vec<Vec<i32>> {
        (0..(1 << nums.len()))
            .map(|mask| {
                nums.iter()
                    .enumerate()
                    .filter_map(|(i, &v)| if mask & (1 << i) != 0 { Some(v) } else { None })
                    .collect()
            })
            .collect()
    }

    pub fn get_row(row_index: i32) -> Vec<i32> {
        Self::gen_pascal_triangle().nth(row_index as usize).unwrap()
    }

    pub fn generate(num_rows: i32) -> Vec<Vec<i32>> {
        Self::gen_pascal_triangle()
            .take(num_rows as usize)
            .collect()
    }

    // 杨辉三角迭代器
    fn gen_pascal_triangle() -> impl Iterator<Item = Vec<i32>> {
        std::iter::successors(Some(vec![1]), |pre_row| {
            vec![1]
                .into_iter()
                .chain(pre_row.windows(2).map(|v| v.into_iter().sum()))
                .chain(vec![1])
                .collect::<Vec<i32>>()
                .into()
        })
    }

    pub fn fib(n: i32) -> i32 {
        assert!(n >= 0 && n <= 100);
        (0..n)
            .fold((0u128, 1), |(a, b), _| (b, (a + b) % 1000000007))
            .0 as i32
    }

    pub fn min_path_sum(grid: Vec<Vec<i32>>) -> i32 {
        grid.iter()
            .fold(None::<Vec<i32>>, |last_row, row| {
                row.into_iter()
                    .enumerate()
                    .scan(None, |left, (i, &v)| {
                        *left = match (&left, &last_row) {
                            (None, None) => v,                         // first cell
                            (Some(left), None) => left + v,            // first row
                            (None, Some(last_row)) => last_row[i] + v, // first column
                            (Some(left), Some(last_row)) => std::cmp::min(*left, last_row[i]) + v,
                        }
                        .into();
                        *left
                    })
                    .collect::<Vec<_>>()
                    .into()
            })
            .unwrap()
            .pop()
            .unwrap()
    }

    pub fn is_subsequence(s: String, t: String) -> bool {
        t.chars()
            .try_fold(s.chars().peekable(), |mut s, t| match s.peek() {
                Some(&c) if c == t => {
                    s.next();
                    Some(s)
                }
                Some(_) => Some(s),
                _ => None,
            })
            .map_or(true, |mut s| s.next().is_none())
    }

    pub fn num_identical_pairs(nums: Vec<i32>) -> i32 {
        nums.into_iter()
            .fold(std::collections::HashMap::new(), |mut mp, v| {
                mp.entry(v).and_modify(|x| *x += 1).or_insert(1);
                mp
            })
            .values()
            .map(|v| v * (v - 1) / 2)
            .sum()
    }

    pub fn is_toeplitz_matrix(matrix: Vec<Vec<i32>>) -> bool {
        matrix
            .windows(2)
            .all(|w| w[0].iter().zip(w[1].iter().skip(1)).all(|(a, b)| a == b))
    }

    pub fn is_monotonic(nums: Vec<i32>) -> bool {
        let (mut inc, mut dec) = (false, false);
        nums.windows(2).all(|w| {
            inc |= w[0] < w[1];
            dec |= w[0] > w[1];
            !(inc & dec)
        })
    }

    pub fn remove_duplicates(s: String) -> String {
        let mut stack = vec![];
        for b in s.into_bytes() {
            if stack.last().map_or(false, |&c| c == b) {
                stack.pop();
            } else {
                stack.push(b);
            }
        }
        unsafe { String::from_utf8_unchecked(stack) }
    }
    pub fn search_matrix(matrix: Vec<Vec<i32>>, target: i32) -> bool {
        matrix
            .binary_search_by(|row| row[0].cmp(&target))
            .map_or_else(
                |less_index| {
                    less_index != 0 && matrix[less_index - 1].binary_search(&target).is_ok()
                },
                |_| true,
            )
    }

    pub fn top_k_frequent(words: Vec<String>, k: i32) -> Vec<String> {
        use std::cmp::Reverse;
        use std::collections::{BinaryHeap, HashMap};
        let mut mp = HashMap::<String, usize>::new();
        words
            .into_iter()
            .for_each(|w| *mp.entry(w).or_default() += 1);

        mp.into_iter()
            .map(|(w, n)| (Reverse(n), w))
            .collect::<BinaryHeap<_>>()
            .into_sorted_vec()
            .into_iter()
            .map(|(_, w)| w)
            .take(k as usize)
            .collect()
    }

    pub fn group_anagrams(strs: Vec<String>) -> Vec<Vec<String>> {
        use std::collections::HashMap;
        let mut mp = HashMap::<[i32; 26], Vec<String>>::new();
        for s in strs {
            let mut counter = [0i32; 26];
            for b in s.bytes() {
                counter[(b - b'a') as usize] += 1;
            }
            mp.entry(counter).or_default().push(s);
        }
        mp.into_iter().map(|(_, v)| v).collect()
    }

    pub fn array_strings_are_equal(word1: Vec<String>, word2: Vec<String>) -> bool {
        word1
            .into_iter()
            .flat_map(|s| s.into_bytes())
            .eq(word2.into_iter().flat_map(|s| s.into_bytes()))
    }
    // 6. Z 字形变换
    // https://leetcode.cn/problems/zigzag-conversion/description/
    pub fn convert(s: String, num_rows: i32) -> String {
        let num_rows = num_rows as usize;
        let mut rows = vec![String::new(); num_rows];
        // z字形往复的迭代器，01232101232......
        let iter = (0..num_rows).chain((1..num_rows - 1).rev()).cycle();
        iter.zip(s.chars()).for_each(|(i, c)| rows[i].push(c));

        rows.into_iter().collect()
    }
}
