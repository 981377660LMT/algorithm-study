static mut MY_STATIC: i32 = 42;

fn main() {
    // 2-4 元组与数组
    {
        let tuple = (1, "hello", 3.4);
        println!("x.0 = {}, {}", tuple.0, tuple.1);

        let mut mutable_tuple = (1, 2);
        mutable_tuple.0 = 3;
        println!("x.0 = {}", mutable_tuple.0);

        let mut arr = [1, 2, 3, 4, 5];
        arr[0] = 99;
        for i in 0..arr.len() {
            println!("arr[{}] = {}", i, arr[i]);
        }

        let num = [2; 3];

        let s = String::from("hello");
    }
}
