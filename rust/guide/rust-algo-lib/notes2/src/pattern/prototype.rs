// 原型是一种创建型设计模式， 使你能够`复制`对象，
// 甚至是复杂对象， 而又无需使代码依赖它们所属的类
// !让你能够复制已有对象， 而又无需使代码依赖它们所属的类
//
// Rust 具有内置 std::clone::Clone 特性，
// 具有许多适用于各种类型的实现（通过 #[derive(Clone)] ）。
// !因此，Prototype 模式开箱即用。

#[derive(Clone, Debug)]
struct Circle {
    pub x: i32,
    pub y: i32,
    pub radius: u32,
}

fn main() {
    let circle1 = Circle {
        x: 10,
        y: 20,
        radius: 30,
    };

    let mut circle2 = circle1.clone();
    circle2.x = 100;

    println!("{:#?}", circle1);
    println!("{:#?}", circle2);
}
