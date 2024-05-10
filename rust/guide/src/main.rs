mod error_handling;
mod lifecycle;
mod misc;
mod pub_use;
mod smart_pointer;
mod test_closoure;
mod test_collections;
mod test_concurrent;
mod test_func;
mod test_iterator;
mod web_server;

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
        let s1: String = String::from("hello");
        let s2: String = s1; // s1 的所有权转移给s2
                             // println!("{s1}"); //    value borrowed here after move
        let s3: String = s2.clone();
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

        fn first_word(s: &str) -> &str {
            let bytes = s.as_bytes();
            for (i, &item) in bytes.iter().enumerate() {
                if item == b' ' {
                    return &s[0..i];
                }
            }
            &s[..]
        }

        let back = first_word("hello world");
        println!("first word: {}", back);
    }

    // 3.2 String 和 &str
    {
        let name = String::from("hello");
        let course = "Rust".to_owned();
        let new_name = name.replace("hello", "world");
        println!("new_name = {}", new_name);
        let rust = "Rust";

        struct User {
            name: String,
            age: u32,
        }

        let user = User {
            name: "Tom".to_string(),
            age: 18,
        };
        println!("user:{}", { user.name });
    }
}

fn first_word_index(s: &str) -> &str {
    let bytes = s.as_bytes();
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    let d = Rectangle {
        width: 1,
        height: 2,
    };

    &s[..]
}

fn demo() {
    let s = Rectangle::new(1, 1);
}

#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn new(width: u32, height: u32) -> Rectangle {
        Rectangle { width, height }
    }

    fn area(&self) -> u32 {
        self.width * self.height
    }

    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }
}

enum IpAddrKind {
    V4,
    V6,
}

fn test_option() {
    let some_num: Option<i32> = Some(5);
    let absent_num: Option<i32> = None;
}

enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
}

fn test_match(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => {
            println!("Lucky penny!");
            1
        }
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}
