// 8. 字符串转换整数 (atoi)
// https://leetcode.cn/problems/string-to-integer-atoi/description/?envType=study-plan-v2&envId=programming-skills
// https://leetcode.cn/problems/string-to-integer-atoi/solutions/522478/rust-you-xian-zhuang-tai-ji-iterator-by-wv2w1/?envType=study-plan-v2&envId=programming-skills
// 请你来实现一个 myAtoi(string s) 函数，使其能将字符串转换成一个 32 位有符号整数。
// 函数 myAtoi(string s) 的算法如下：
// 空格：读入字符串并丢弃无用的前导空格（" "）
// 符号：检查下一个字符（假设还未到字符末尾）为 '-' 还是 '+'。如果两者都不存在，则假定结果为正。
// 转换：通过跳过前置零来读取该整数，直到遇到非数字字符或到达字符串的结尾。如果没有读取数字，则结果为0。
// 舍入：如果整数数超过 32 位有符号整数范围 [−231,  231 − 1] ，需要截断这个整数，使其保持在这个范围内。具体来说，小于 −231 的整数应该被舍入为 −231 ，大于 231 − 1 的整数应该被舍入为 231 − 1 。
// 返回整数作为最终结果。

// Pattern Matching 语法写状态机非常方便，尽管还不如 Erlang/Elixir 的方便。
// Rust 的 enum 类型用来实现状态机的状态，把状态转移实现在 enum 类型内部，通过函数接口暴露给外部。
// 使用 Iterator 的 combinator 处理高阶函数。

// 需要注意的是在解析成数字后的运算需要检查是否 overflow，Rust 标准库本身就实现了这些安全的检查操作，结合 Option, Result 的 combinator 操作。
//
// Option.map(f) -> 将 Option 中的值应用 f 函数，返回 Option
// Option.map_err(f) -> 将 Result 中的 Err 应用 f 函数，返回 Result
// Option.map_or(default, f) -> 将 Option 中的值应用 f 函数，如果 Option 为 None 则返回 default
// Option.map_or_else(default, f) -> 将 Option 中的值应用 f 函数，如果 Option 为 None 则返回 default 函数的值

struct Solution;

enum State {
    Start,
    Sign(i32),
    Num(i32),
}

impl State {
    fn new() -> Self {
        State::Start
    }

    fn into_result(self) -> i32 {
        match self {
            State::Num(n) => n,
            _ => 0,
        }
    }

    fn next(self, c: char) -> Result<Self, i32> {
        match self {
            State::Start => match c {
                ' ' => Ok(State::Start),
                '+' => Ok(State::Sign(1)),
                '-' => Ok(State::Sign(-1)),
                c @ '0'..='9' => Ok(State::Num(c.to_digit(10).unwrap() as i32)),
                _ => Err(0),
            },
            State::Sign(sign) => match c {
                '0' => Ok(State::Sign(sign)),
                c @ '1'..='9' => Ok(State::Num(sign * c.to_digit(10).unwrap() as i32)),
                _ => Err(0),
            },
            State::Num(n) => match c {
                c @ '0'..='9' => {
                    if n >= 0 {
                        n.checked_mul(10)
                            .and_then(|v| v.checked_add(c.to_digit(10).unwrap() as i32))
                            .map(State::Num)
                            .ok_or(i32::MAX)
                    } else {
                        n.checked_mul(10)
                            .and_then(|v| v.checked_sub(c.to_digit(10).unwrap() as i32))
                            .map(State::Num)
                            .ok_or(i32::MIN)
                    }
                }
                _ => Err(n),
            },
        }
    }
}

impl Solution {
    pub fn my_atoi(s: String) -> i32 {
        s.chars()
            .try_fold(State::new(), |state, c| state.next(c))
            .map_or_else(|n| n, |state| state.into_result())
    }
}
