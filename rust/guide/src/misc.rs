fn main() {
    let slice = vec![1, 2, 3, 4];
    let foo = &slice[2..3];
    let (a, _) = slice.split_at(2);
    println!("{:?}", a); // [1, 2]

    let mut x = 42;
    let x_ref1 = &x;
    let x_ref2 = &mut x;
    // error: cannot borrow `x` as mutable because it is also borrowed as immutable
    println!("x_ref1 = {}", x_ref1);
}
