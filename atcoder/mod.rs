use crate::rust::{ds::union_find::UnionFind, misc::io::run_with_io};
use std::io::Write;

fn main() {
    let mut uf = UnionFind::new(10);
    println!("{}", uf.is_connected(0, 1));
    run_with_io(|reader, writer| {
        let a: i32 = reader.next();
        writeln!(writer, "{}", a).unwrap();
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        main()
    }
}
