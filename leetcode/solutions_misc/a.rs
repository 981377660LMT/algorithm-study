use std::collections::{HashMap, HashSet};
use std::iter::{Fuse, Peekable};

type Result<T> = std::result::Result<T, String>;

impl Solution {
    pub fn evaluate(expression: String) -> i32 {
        let mut e = Evaluation::new(tokenize(&expression));
        // assume all inputs will be valid
        let res = e.eval(None).unwrap_or_else(|e| panic!("{}", e));
        // 如果还有符号，说明输入是非法的
        assert!(e.is_end(), "Unexpected input");

        res
    }
}

struct Evaluation<T: Iterator> {
    /// stacks for each level local variables
    stacks: Vec<HashMap<String, i32>>,

    /// entire namespace in the current level
    vars: HashMap<String, i32>,

    tokens: Peekable<Fuse<T>>,
}

impl<T: Iterator> Evaluation<T> {
    fn new(tokens: T) -> Self {
        Self {
            stacks: Vec::new(),
            vars: HashMap::new(),
            tokens: tokens.fuse().peekable(),
        }
    }
}

impl<T: Iterator<Item = Token>> Evaluation<T> {
    #[inline]
    fn expect(&mut self, expect: Token) -> Result<()> {
        self.tokens
            .next()
            .ok_or_else(|| format!("Invalid, missing required token: {:?}", expect))
            .and_then(|token| {
                if token == expect {
                    Ok(())
                } else {
                    Err(format!(
                        "Invalid, expect token: {:?}, found: {:?}",
                        expect, token
                    ))
                }
            })
    }

    #[inline]
    fn is_end(&mut self) -> bool {
        self.tokens.peek().is_none()
    }

    // Evaluate expression
    fn eval(&mut self, local: Option<&HashSet<String>>) -> Result<i32> {
        match self.tokens.next() {
            // 遇到左括号，递归进入
            // 递归返回的时候，这一层的表达式求值完成，需要把右括号也弹出
            Some(Token::LeftParenthesis) => {
                self.eval(local).and_then(|res| {
                    // expect RightParenthesis
                    self.expect(Token::RightParenthesis)?;
                    Ok(res)
                })
            }

            // 处理 `let` 语句
            Some(Token::Let) => {
                // 进入 let 语句的时候，需要入栈
                // 退出 let 语句的时候，需要出栈
                self.push_stack(local);
                let mut local = HashSet::new();
                let res = self.eval_let(&mut local)?;
                self.pop_stack(Some(&local));
                Ok(res)
            }

            // `add` 语句，有两个参数，所以可以直接递归求值
            Some(Token::Add) => self
                .eval(local)
                .and_then(|e1| self.eval(local).map(|e2| e1 + e2)),

            // `mult` 语句，同 `add`
            Some(Token::Mult) => self
                .eval(local)
                .and_then(|e1| self.eval(local).map(|e2| e1 * e2)),

            // 常量表达式，直接求值返回
            Some(Token::Const(num)) => Ok(num),

            // 标识符，即变量表达式，直接求值返回
            Some(Token::Ident(var)) => self
                .vars
                .get(&var)
                .map(|&v| v)
                .ok_or_else(|| format!("Not found var: `{}`", var)),

            // 其它情况不支持
            _ => Err("Invalid token".into()),
        }
    }

    /// Evaluate `let` expression, from the first item after `let`
    fn eval_let(&mut self, local: &mut HashSet<String>) -> Result<i32> {
        loop {
            // 分两种情况，先处理键值对的赋值语句，然后处理 `let` 语句中最后的 expr 部分
            match self.tokens.peek() {
                // 如果下一个表达式是标识符，那么有以下几种符合的情况
                Some(Token::Ident(_var)) => {
                    let ident = match self.tokens.next().unwrap() {
                        Token::Ident(s) => s,
                        _ => unreachable!(),
                    };

                    match self.tokens.peek() {
                        // assignment a var with a parenthesis expression (...) or a constant or a variable
                        // 赋值表达式，给之前的 ident 赋值，可以是 括号表达式、常量表达式、变量表达式，求值后更新变量空间，同时记录 local 变量名，循环处理
                        Some(Token::LeftParenthesis)
                        | Some(Token::Const(_))
                        | Some(Token::Ident(_)) => {
                            let res = self.eval(Some(&local))?;
                            self.vars.insert(ident.clone(), res);
                            local.insert(ident);
                        }
                        // evaluate the last expression of let
                        // 如果标识符之后接着是右括号，说明这个标识符是 `let` 语句的最后一部分，也就是说这个标识符是变量表达式，可以直接求值返回
                        Some(Token::RightParenthesis) => {
                            return self
                                .vars
                                .get(&ident)
                                .map(|&v| v)
                                .ok_or_else(|| format!("Not found var: `{}`", ident));
                        }
                        _ => {
                            return Err("Invalid let, expect expression".into());
                        }
                    }
                }

                // let expr part is const or is (...) expr
                // 如果最后一个表达式是常量表达式或者是括号表达式，可以直接求值回返，紧接着之后一定是一个右括号，否则就不合法
                Some(Token::Const(_)) | Some(Token::LeftParenthesis) => {
                    return self.eval(Some(&local));
                    // next should be Token::RightParenthesis
                }

                _ => {
                    return Err("Invalid `let` expression".into());
                }
            }
        }
    }

    fn push_stack(&mut self, local: Option<&HashSet<String>>) {
        local.map(|local| {
            let mut stack = HashMap::new();
            for k in local.iter() {
                if let Some(v) = self.vars.get(k) {
                    stack.insert(k.clone(), v.clone());
                }
            }
            self.stacks.push(stack);
        });
    }

    fn pop_stack(&mut self, local: Option<&HashSet<String>>) {
        local.map(|local| {
            for k in local.iter() {
                self.vars.remove(k);
            }

            for stack in &self.stacks {
                for (k, v) in stack {
                    self.vars.insert(k.clone(), v.clone());
                }
            }
            self.stacks.pop();
        });
    }
}

#[derive(Debug, PartialEq)]
enum Token {
    /// 常量
    Const(i32),

    /// 标识符，可以是用于赋值左侧的变量名，也可以是一个求值表达式
    Ident(String),

    /// Keyword `let`
    Let,

    /// Keyword `add`
    Add,

    /// Keyword `mult`
    Mult,

    /// Left parenthesis `(`
    LeftParenthesis,

    /// Right parenthesis `)`
    RightParenthesis,
}

impl Token {
    fn from_str<S: AsRef<str>>(s: S) -> Self {
        let s = s.as_ref();
        match s {
            "(" => Self::LeftParenthesis,
            ")" => Self::RightParenthesis,
            "let" => Self::Let,
            "add" => Self::Add,
            "mult" => Self::Mult,
            c if c.len() > 0 => {
                match c.chars().next().unwrap() {
                    'a'..='z' => Self::Ident(s.into()),
                    _ => Self::Const(s.parse().unwrap()), // assume all expr valid
                }
            }
            _ => {
                panic!("Invalid");
            }
        }
    }
}

fn tokenize<'a>(s: &'a str) -> impl Iterator<Item = Token> + 'a {
    let mut chars = s.chars();
    let mut cache = String::new();
    std::iter::from_fn(move || {
        while let Some(c) = chars.next() {
            match c {
                '(' => return Some(Token::from_str("(")),
                v @ ')' => {
                    if cache.len() > 0 {
                        let next = Some(Token::from_str(cache.drain(..).collect::<String>()));
                        cache.push(v);
                        return next;
                    } else {
                        return Some(Token::from_str(")"));
                    }
                }
                ' ' => {
                    return Some(Token::from_str(cache.drain(..).collect::<String>()));
                }
                c => cache.push(c),
            }
        }

        if cache.len() > 0 {
            Some(Token::from_str(cache.drain(..).collect::<String>()))
        } else {
            None
        }
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        let cases = vec![
            (3, "(add 1 2)"),
            (15, "(mult 3 (add 2 3))"),
            (10, "(let x 2 (mult x 5))"),
            (14, "(let x 2 (mult x (let x 3 y 4 (add x y))))"),
            (2, "(let x 3 x 2 x)"),
            (5, "(let x 1 y 2 x (add x y) (add x y))"),
            (6, "(let x 2 (add (let x 3 (let x 4 x)) x))"),
            (4, "(let a1 3 b2 (add a1 1) b2)"),
            (14, "(let x 2 (mult (let x 3 y 4 (add x y)) x))"),
            (-128534112, "(let x0 -4 x1 2 x2 -4 x3 3 x4 2 x5 3 x6 2 x7 2 x8 -1 x9 -1 (mult (mult (mult x2 -8) (add -5 (let x0 1 x5 -3 (add (add x7 (add (let x0 -5 x9 -4 (add (mult 1 1) -10)) (mult -8 (mult x3 -5)))) (add (let x0 3 x8 -1 (let x0 -1 x9 1 (add x4 -6))) x9))))) (mult (add (mult (add (mult -6 (mult (add x1 x4) -4)) (let x0 -2 x7 4 (mult (mult (let x0 -3 (mult 1 1)) (add (mult 1 1) (mult 1 1))) (mult -5 (mult -9 (mult 1 1)))))) -10) x5) (mult (mult x5 -7) x8))))"),
        ];
        for (expect, arg) in cases {
            assert_eq!(expect, Solution::evaluate(arg.into()));
        }
    }
}
