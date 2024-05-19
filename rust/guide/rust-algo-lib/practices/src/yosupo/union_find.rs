// https://judge.yosupo.jp/problem/unionfind

use std::io::Write;
use templates::{ds::union_find::UnionFind, misc::io::run_with_io};

#[allow(unused)]
fn main() {
    run_with_io(|sc, out| {
        let n: usize = sc.next();
        let q: usize = sc.next();
        let mut uf = UnionFind::new(n);
        for _ in 0..q {
            let t: usize = sc.next();
            let u: usize = sc.next();
            let v: usize = sc.next();
            if t == 0 {
                uf.union(u, v, None);
            } else {
                writeln!(out, "{}", if uf.is_connected(u, v) { 1 } else { 0 }).ok();
            }
        }
    });
}
