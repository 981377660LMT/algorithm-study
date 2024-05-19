mod atcoder;
mod yosupo;

fn main() {
    println!("Hello, world!");
    // 没有返回值的函数在 Rust 中是有单独的定义的：发散函数( diverge function )
    // 用 () 作为 map 的值，表示我们不关注具体的值，只关注 key。
    // 这种用法和 Go 语言的 struct{} 类似，可以作为一个值用来占位，但是完全不占用任何内存。
    let _a = ();
}
