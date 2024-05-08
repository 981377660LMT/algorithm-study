fn main() {
    let s1: String = String::from("hello");
    {
        let s2: String = String::from("world");
        let result: &str = longer(s1.as_str(), s2.as_str());
        println!("The longer string is {}", result);
    }
}

// this function's return type contains a borrowed value, but the signature does not say whether it is borrowed from `s1` or `s2`
// !取生命周期最短的那个
fn longer<'a>(s1: &'a str, s2: &'a str) -> &'a str {
    if s1.len() > s2.len() {
        s1
    } else {
        s2
    }
}
