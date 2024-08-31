struct Solution;

// 编程基础 0 到 1
// https://leetcode.cn/studyplan/programming-skills/
impl Solution {
    pub fn merge_alternately(word1: String, word2: String) -> String {
        let (mut iter1, mut iter2) = (word1.chars().peekable(), word2.chars().peekable());
        let mut res = "".to_string();
        while iter1.peek().is_some() || iter2.peek().is_some() {
            if let Some(c) = iter1.next() {
                res.push(c);
            }
            if let Some(c) = iter2.next() {
                res.push(c);
            }
        }
        res
    }

    pub fn find_the_difference(s: String, t: String) -> char {
        s.as_bytes()
            .iter()
            .chain(t.as_bytes().iter())
            .fold(0, |pre, cur| pre ^ cur) as char
    }

    pub fn str_str(haystack: String, needle: String) -> i32 {
        return haystack.find(&needle).map_or(-1, |x| x as i32);
    }

    pub fn is_anagram(s: String, t: String) -> bool {
        let mut cnt = [0; 26];
        s.bytes().for_each(|b| cnt[b as usize - 97] += 1);
        t.bytes().for_each(|b| cnt[b as usize - 97] -= 1);
        cnt.iter().all(|&x| x == 0)
    }

    pub fn repeated_substring_pattern(s: String) -> bool {
        return (s.clone() + &s)[1..s.len() * 2 - 1].contains(&s);
    }

    pub fn move_zeroes(nums: &mut Vec<i32>) {
        let n = nums.len();
        nums.retain(|&x| x != 0);
        nums.resize(n, 0);
    }

    pub fn plus_one(mut digits: Vec<i32>) -> Vec<i32> {
        let mut carry = false;
        for i in (0..digits.len()).rev() {
            match digits[i] {
                9 => {
                    digits[i] = 0;
                    carry = true;
                }
                _ => {
                    digits[i] += 1;
                    carry = false;
                    break;
                }
            }
        }
        if carry {
            digits.insert(0, 1);
        }
        digits
    }

    pub fn array_sign(nums: Vec<i32>) -> i32 {
        if nums.contains(&0) {
            0
        } else {
            nums.iter()
                .fold(1, |pre, &cur| if cur < 0 { -pre } else { pre })
        }
    }

    pub fn can_make_arithmetic_progression(mut arr: Vec<i32>) -> bool {
        arr.sort_unstable();
        arr.windows(3).all(|w| w[1] + w[1] == w[0] + w[2])
    }

    pub fn is_monotonic(nums: Vec<i32>) -> bool {
        return nums.iter().zip(nums.iter().skip(1)).all(|(a, b)| a <= b)
            || nums.iter().zip(nums.iter().skip(1)).all(|(a, b)| a >= b);
    }

    pub fn roman_to_int(s: String) -> i32 {
        s.bytes()
            .rev()
            .fold((0, 0), |res, cur| {
                let n = match cur {
                    b'I' => 1,
                    b'V' => 5,
                    b'X' => 10,
                    b'L' => 50,
                    b'C' => 100,
                    b'D' => 500,
                    b'M' => 1000,
                    _ => -9999,
                };
                (if n < res.1 { res.0 - n } else { res.0 + n }, n)
            })
            .0
    }

    pub fn length_of_last_word(s: String) -> i32 {
        s.into_bytes()
            .into_iter()
            .rev()
            .skip_while(|&c| c == b' ')
            .take_while(|&c| c != b' ')
            .count() as i32
    }

    pub fn to_lower_case(s: String) -> String {
        s.to_ascii_lowercase()
    }

    pub fn cal_points(operations: Vec<String>) -> i32 {
        let mut stack = vec![];
        for op in operations {
            match op.as_bytes()[0] {
                b'C' => {
                    stack.pop();
                }
                b'D' => {
                    stack.push(stack.last().unwrap() * 2);
                }
                b'+' => {
                    stack.push(stack[stack.len() - 1] + stack[stack.len() - 2]);
                }
                _ => {
                    stack.push(op.parse::<i32>().unwrap());
                }
            }
        }
        stack.into_iter().sum()
    }

    pub fn judge_circle(moves: String) -> bool {
        moves.chars().fold((0, 0), |(x, y), c| match c {
            'U' => (x, y + 1),
            'D' => (x, y - 1),
            'L' => (x - 1, y),
            'R' => (x + 1, y),
            _ => (x, y),
        }) == (0, 0)
    }

    pub fn maximum_wealth(accounts: Vec<Vec<i32>>) -> i32 {
        return accounts
            .into_iter()
            .map(|v| v.into_iter().sum())
            .max()
            .unwrap();
    }

    pub fn diagonal_sum(mat: Vec<Vec<i32>>) -> i32 {
        let n = mat.len();
        (0..n).fold(0, |sum, i| sum + mat[i][i] + mat[i][n - i - 1])
            - (n & 1) as i32 * mat[n / 2][n / 2]
    }

    pub fn lemonade_change(bills: Vec<i32>) -> bool {
        let (mut a, mut b) = (0, 0);
        for bill in bills {
            match bill {
                5 => a += 1,
                10 => {
                    a -= 1;
                    b += 1;
                }
                _ => {
                    if b > 0 {
                        b -= 1;
                        a -= 1;
                    } else {
                        a -= 3;
                    }
                }
            }
            if a < 0 {
                return false;
            }
        }
        true
    }

    pub fn check_straight_line(coordinates: Vec<Vec<i32>>) -> bool {
        let (x0, y0) = (coordinates[0][0], coordinates[0][1]);
        let (x1, y1) = (coordinates[1][0], coordinates[1][1]);
        let (dy, dx) = (y1 - y0, x1 - x0);
        coordinates
            .into_iter()
            .all(|v| dy * (v[0] - x0) == dx * (v[1] - y0))
    }

    pub fn my_pow(x: f64, n: i32) -> f64 {
        x.powf(n as f64)
    }
}
