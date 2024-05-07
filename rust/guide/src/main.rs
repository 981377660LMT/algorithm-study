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

    // 3-1 内存管理模型
    {
        // copy/move

        // copy 只对基础类型有效
        let c1 = 1;
        let c2 = c1;
        println!("c1 = {}, c2 = {}", c1, c2);

        // move 赋值是所有权转移，原来的变量就不能再使用
        // 如果需要拷贝，需要.clone()方法复制
        let s1 = String::from("hello");
        let s2 = s1;
        // println!("{s1}"); //    value borrowed here after move
        let s3 = s2.clone();
        println!("s2 = {}, s3 = {}", s2, s3);
        // take_ownership(s3); // s3的所有权转移到函数内
        // println!("{s3}"); //    value borrowed here after move
        println!("s3 length = {}", get_length(&s3));

        fn take_ownership(s: String) {
            // 函数结束后，s的内存会被释放
            println!("take_ownership: {}", s);
        }

        fn get_length(s: &String) -> u32 {
            s.len() as u32
        }

        // fn dangle() -> &String {
        //   let s = String::from("hello");
        //   cannot return reference to local variable `s`
        //   returns a reference to data owned by the current function
        //   &s
        // }
    }
}
