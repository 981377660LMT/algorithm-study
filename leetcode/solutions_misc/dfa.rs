// 65. 有效数字(dfa解法)
// https://leetcode.cn/problems/valid-number/solutions/2078708/rust-dfa-100-100-by-selenium34-h1ss/?envType=study-plan-v2&envId=programming-skills
// DFA 五元组
// 状态集合 Q: enum State
// 输入集合 Σ: enum Input
// 转移函数 Q×Σ→Q: State::next(&self, next: Input) -> Self
// 起始状态 q0: State::Start
// 结束状态集合 F: State::Integer | State::Dot | State::Decimal | State::EInteger

// From 和 Into
// https://rustwiki.org/zh-CN/rust-by-example/conversion/from_into.html

struct Solution;

enum State {
    Start,
    Sign,
    Integer,
    Dot,
    InitialDot,
    Decimal,
    E,
    ESign,
    EInteger,
    Illegal,
}

enum Input {
    Digit,
    Sign,
    Dot,
    E,
    Other,
}

impl State {
    fn new() -> Self {
        State::Start
    }

    fn next(self, next: Input) -> Self {
        match self {
            State::Start => match next {
                Input::Digit => State::Integer,
                Input::Dot => State::InitialDot,
                Input::Sign => State::Sign,
                _ => State::Illegal,
            },
            State::Sign => match next {
                Input::Digit => State::Integer,
                Input::Dot => State::InitialDot,
                _ => State::Illegal,
            },
            State::Integer => match next {
                Input::Digit => State::Integer,
                Input::Dot => State::Dot,
                Input::E => State::E,
                _ => State::Illegal,
            },
            State::Dot => match next {
                Input::Digit => State::Decimal,
                Input::E => State::E,
                _ => State::Illegal,
            },
            State::InitialDot => match next {
                Input::Digit => State::Decimal,
                _ => State::Illegal,
            },
            State::Decimal => match next {
                Input::Digit => State::Decimal,
                Input::E => State::E,
                _ => State::Illegal,
            },
            State::E => match next {
                Input::Digit => State::EInteger,
                Input::Sign => State::ESign,
                _ => State::Illegal,
            },
            State::ESign => match next {
                Input::Digit => State::EInteger,
                _ => State::Illegal,
            },
            State::EInteger => match next {
                Input::Digit => State::EInteger,
                _ => State::Illegal,
            },
            _ => State::Illegal,
        }
    }

    fn accept(self) -> bool {
        match self {
            State::Integer | State::Dot | State::Decimal | State::EInteger => true,
            _ => false,
        }
    }
}

impl From<char> for Input {
    fn from(c: char) -> Self {
        match c {
            '0'..='9' => Input::Digit,
            '+' | '-' => Input::Sign,
            '.' => Input::Dot,
            'e' | 'E' => Input::E,
            _ => Input::Other,
        }
    }
}

impl Solution {
    pub fn is_number(s: String) -> bool {
        s.chars()
            .fold(State::new(), |state, c| state.next(c.into()))
            .accept()
    }
}
