use std::{
    fs::File,
    io::{self, ErrorKind, Read},
};

fn main() {
    handle_error();
}

fn handle_error() {
    let f: Result<File, std::io::Error> = File::open("hello.txt");
    // let f: File = match f {
    //     Ok(file) => file,
    //     Err(error) => {
    //         panic!("Problem opening the file: {:?}", error)
    //     }
    // };

    // let f = File::open("hello.txt").expect("Failed to open hello.txt");

    let f: File = File::open("hello.txt").unwrap_or_else(|err: std::io::Error| {
        if err.kind() == ErrorKind::NotFound {
            File::create("hello.txt").unwrap_or_else(|err| {
                panic!("Problem creating the file: {:?}", err);
            })
        } else {
            panic!("Problem opening the file: {:?}", err);
        }
    });
}

fn read_username_from_file() -> Result<String, io::Error> {
    // ? 运算符只能用于返回Result的函数
    // 表示如果Result是Ok，则返回Ok的值，否则返回Err的值
    let mut s = String::new();
    File::open("foo.text")?.read_to_string(&mut s)?;
    Ok(s)
}
