use std::slice;

fn main() {
    let x;
    x = 22;
    let x = 33;
    let pair: (char, i32) = ('a', 1);
    let x = vec![1, 2, 3, 4, 5]
        .iter()
        .map(|v| v + 1)
        .fold(0, |x, y| x + y);
    println!("x = {}", x);
}
