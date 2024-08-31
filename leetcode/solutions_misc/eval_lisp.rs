// 736. Lisp 语法解析
// !先 parse 成 AST 再 eval
// https://leetcode.cn/problems/parse-lisp-expression/solutions/600836/rustxian-parse-cheng-ast-zai-eval-by-ant-qbz2/
// !这样写虽然比读入token的同时eval的做法代码长一点，但是更不容易出错也更好扩展

// option.copied() -> Option<&T> -> Option<T>

use std::{collections::HashMap, f32::consts::E};

/// 变量绑定
type Bindings = HashMap<String, i32>;

/// 记录变量绑定的环境
struct Env(Vec<Bindings>);

impl Env {
    fn new() -> Self {
        Self(vec![])
    }

    fn insert(&mut self, k: String, v: i32) {
        self.0.last_mut().and_then(|bindings| bindings.insert(k, v));
    }

    fn find(&self, var: &str) -> Option<i32> {
        self.0
            .iter()
            .rev()
            .find_map(|bindings| bindings.get(var).copied())
    }

    fn push_frame(&mut self) {
        self.0.push(HashMap::new());
    }

    fn pop_frame(&mut self) -> Option<Bindings> {
        self.0.pop()
    }
}

/// 语法树, 描述Lisp的基本结构
enum Expr {
    Num(i32),
    Op(String),
    List(Vec<Expr>),
}

impl Expr {
    /// 表达式解析
    fn evaluate(&self, env: &mut Env) -> i32 {
        use Expr::*;

        match self {
            Num(n) => *n,
            // 对于操作符，在栈中查找对应的变量
            Op(s) => env.find(s).unwrap(),
            // 对于列表，递归求值
            List(list) => Self::eval_list(list, env),
        }
    }

    fn eval_list(list: &[Expr], env: &mut Env) -> i32 {
        use Expr::*;

        // let 形式的求值
        let eval_let = |l: &[Expr], env: &mut Env| -> i32 {
            env.push_frame();

            // 变量和值成对，单出来的是最终需要求值的表达式
            let pairs = l.chunks_exact(2);
            let tail_expr = &pairs.remainder()[0];

            // 对于每对变量和值，对值表达式求值后加入环境
            pairs.for_each(|pair| {
                if let [Op(k), expr] = pair {
                    let v = expr.evaluate(env);
                    env.insert(k.clone(), v);
                }
            });

            let res = tail_expr.evaluate(env);
            env.pop_frame();
            res
        };

        match list {
            // let 形式
            [Op(s), l @ ..] if s == "let" => eval_let(l, env),
            [Op(s), a, b] if s == "add" => a.evaluate(env) + b.evaluate(env),
            [Op(s), a, b] if s == "mult" => a.evaluate(env) * b.evaluate(env),
            _ => unreachable!(),
        }
    }
}

/// 将表达式解析成语法树
fn parse(expr: &str) -> Expr {
    // 解析atom，能解析成数字就解析成数字，否则解析成操作符
    let parse_atom = |atom: String| atom.parse().map_or(Expr::Op(atom), Expr::Num);

    let mut list = vec![];
    let mut stacks = vec![];
    let mut atom = String::new();

    for c in expr.chars() {
        match c {
            '(' => {
                stacks.push(list);
                list = vec![];
            }
            ')' => {
                if !atom.is_empty() {
                    list.push(parse_atom(atom));
                    atom = String::new();
                }
                let mut parent = stacks.pop().unwrap();
                parent.push(Expr::List(list));
                list = parent;
            }
            ' ' => {
                if !atom.is_empty() {
                    list.push(parse_atom(atom));
                    atom = String::new();
                }
            }
            // 其他字符作为atom的一部分
            _ => atom.push(c),
        }
    }

    if !atom.is_empty() {
        // 只有一个atom的情况
        parse_atom(atom)
    } else {
        debug_assert!(list.len() == 1);
        list.pop().unwrap()
    }
}

struct Solution;
impl Solution {
    pub fn evaluate(expression: String) -> i32 {
        parse(&expression).evaluate(&mut Env::new())
    }
}
