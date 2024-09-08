fn main() {
    let mut a = vec![1, 2, 3, 4, 5];
    a.retain(|&x| x % 2 == 0);
    println!("{:?}", a); // [2, 4]
}
