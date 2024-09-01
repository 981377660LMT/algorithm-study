fn main() {
    println!("Hello, world!");
    demo();
}

// ---------- begin scannner ----------
#[allow(dead_code)]
mod scanner {
    use std::str::FromStr;
    pub struct Scanner<'a> {
        it: std::str::SplitWhitespace<'a>,
    }
    impl<'a> Scanner<'a> {
        pub fn new(s: &'a str) -> Scanner<'a> {
            Scanner {
                it: s.split_whitespace(),
            }
        }
        pub fn next<T: FromStr>(&mut self) -> T {
            self.it.next().unwrap().parse::<T>().ok().unwrap()
        }
        pub fn next_bytes(&mut self) -> Vec<u8> {
            self.it.next().unwrap().bytes().collect()
        }
        pub fn next_chars(&mut self) -> Vec<char> {
            self.it.next().unwrap().chars().collect()
        }
        pub fn next_vec<T: FromStr>(&mut self, len: usize) -> Vec<T> {
            (0..len).map(|_| self.next()).collect()
        }
    }
}
// ---------- end scannner ----------

pub fn setup_io() -> (
    scanner::Scanner<'static>,
    std::io::BufWriter<std::io::Stdout>,
) {
    use std::io::Read;
    static mut S: &mut String = &mut String::new();
    std::io::stdin().read_to_string(unsafe { S }).unwrap();
    let sc = scanner::Scanner::new(unsafe { S });
    let out = std::io::stdout();
    let out = std::io::BufWriter::new(out);
    (sc, out)
}

#[allow(unused)]
fn demo() {
    use std::io::Write; // !for writeln!

    let (mut input, mut output) = setup_io();
    let n: usize = input.next();
    writeln!(output, "{}", n).ok();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        main();
    }
}
