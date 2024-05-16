fn main() {}

// 关联类型 Associated Types
// 同一特性实现在不同类型上时，可以具有不同的函数签名

trait Foo {
    type BarType;
    fn func(arg: Self::BarType);
}

struct FooImpl1;
struct FooImpl2;

impl Foo for FooImpl1 {
    type BarType = i32;
    fn func(arg: Self::BarType) {
        println!("FooImpl1: {}", arg);
    }
}

impl Foo for FooImpl2 {
    type BarType = String;
    fn func(arg: Self::BarType) {
        println!("FooImpl2: {}", arg);
    }
}

fn test_associated_types() {
    FooImpl1::func(1);
    FooImpl2::func("hello".to_string());
}
