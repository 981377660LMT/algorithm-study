fn main() {
    // test_vector();
    // test_enum()
    // test_string();
    test_hashmap()
}

fn test_vector() {
    let v_with_macro: Vec<i32> = vec![1, 2, 3];
    let b: Option<&i32> = v_with_macro.get(2);
    if let Some(c) = b {
        println!("{}", c);
    }

    let mut v: Vec<i32> = Vec::new();
    v.push(2);
    let first: &i32 = &v[0];
    // v.push(22); // cannot borrow `v` as mutable because it is also borrowed as immutable
    println!("{:?}", first);
}

fn test_enum() {
    #[derive(Debug)]
    enum SpreadSheetCell {
        Int(i32),
        Float(f64),
        Text(String),
    }

    let row = vec![
        SpreadSheetCell::Int(3),
        SpreadSheetCell::Float(3.4),
        SpreadSheetCell::Text(String::from("hello")),
    ];
    println!("{:#?}", row);
}

fn test_string() {
    let mut s = "abc".to_string();
    s.push_str("def");

    let s1 = "abc".to_string();
    let s2 = "def".to_string();
    let s3 = s1 + &s2; // add 方法

    let s4 = format!("{}-{}", s2, s2); // format!宏
    println!("{}", s4);

    // 访问
    // unicode标量值
    let len = String::from("hello").len();
    let len2 = "你好".len();
    println!("len={}, len2={}", len, len2);

    // bytes
    for b in "你好".bytes() {
        println!("{}", b); // 228 189 160 229 165 189
    }
    // scalar value
    for c in "你好".chars() {
        println!("{}", c); // 你 好
    }

    let s = "你好".to_ascii_lowercase();
    //     byte index 1 is not a char boundary; it is inside '你' (bytes 0..3) of `你好`
    // note: run with `RUST_BACKTRACE=1` environment variable to display a backtrace
    // println!("{}", &s[0..1]); // 你
}

fn test_hashmap() {
    use std::collections::HashMap;
    let mut mp: HashMap<String, i32> = HashMap::new();
    mp.insert("hello".to_string(), 1);

    let keys = vec![String::from("blue"), String::from("yellow")];
    let values = vec![10, 20];
    let mut mp: HashMap<_, _> = keys.iter().zip(values.iter()).collect();
    println!("{:?}", mp);

    mp.entry(&"hella".to_string()).or_insert(&30);
}
