// 全局变量在 Rust 中被称为 静态（static）变量
//
// 完全限定语法定义为：
// <Type as Trait>::function(receiver_if_method, next_arg, ...);
//
// !fn 被称为 函数指针（function pointer），不要与闭包 trait 的 Fn 相混淆
// fn 是一个类型而不是一个 trait，所以直接指定 fn 作为参数而不是声明一个带有 Fn 作为 trait bound 的泛型参数

fn main() {
    static a: i32 = 2;
    println!("{}", do_twice(add_one, 5));

    fn foo() {
        println!("{}", a);
    }

    fn add_one(x: i32) -> i32 {
        x + 1
    }

    fn do_twice(f: fn(i32) -> i32, arg: i32) -> i32 {
        f(arg) + f(arg)
    }
}
