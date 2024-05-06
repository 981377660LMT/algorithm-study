static mut MY_STATIC: i32 = 42;

fn main() {
    let a1: i32 = -13;
    let a2: i32 = 0x3f;
    let a3: i32 = 0o13;
    let a4: i32 = 0b1101;
    println!("a1: {}, a2: {}, a3: {}, a4: {}", a1, a2, a3, a4);

    println!("u32 max: {}", std::mem::size_of::<u32>());

    let 微笑: char = '😊';
    println!("微笑: {微笑}");
}
