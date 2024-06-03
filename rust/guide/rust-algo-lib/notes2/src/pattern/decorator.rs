// 由于目标对象和装饰器遵循同一接口， 因此你可用装饰来对对象进行无限次的封装。
// 结果对象将获得所有封装器叠加而来的行为。

// Input streams decoration
// A buffered reader decorates a vector reader adding buffered behavior.
use std::io::{BufReader, Cursor, Read};

fn main() {
    let mut buf = [0u8; 10];

    // A buffered reader decorates a vector reader which wraps input data.
    let mut input = BufReader::new(Cursor::new("Input data"));

    input.read(&mut buf).ok();

    print!("Read from a buffered reader: ");

    for byte in buf {
        print!("{}", char::from(byte));
    }

    println!();
}
