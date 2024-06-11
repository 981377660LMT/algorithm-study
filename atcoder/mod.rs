use templates::{
    ds::union_find::UnionFind,
    misc::io::{run_with_io, scanner::Scanner},
};

use std::io::{BufWriter, Stdout, Write};

fn main() {
    run_with_io(|reader, writer| {
        run(reader, writer);
    });
}

fn run(reader: &mut Scanner, writer: &mut BufWriter<Stdout>) {
    let a: i32 = reader.next();
    let b: i32 = reader.next();
    writeln!(writer, "{}", a).unwrap();
    writeln!(writer, "{}", b).unwrap();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        main();
    }
}
