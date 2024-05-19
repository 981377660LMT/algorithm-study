use proconio::input;
use templates::ds::union_find::UnionFind;

#[proconio::fastout]
fn main() {
    input! {
      n: usize,
      q: usize,
    }
    let mut uf = UnionFind::new(n);
    for _ in 0..q {
        input! {
          t: usize,
          u: usize,
          v: usize,
        }
        if t == 0 {
            uf.union(u, v, None);
        } else {
            println!("{}", if uf.is_connected(u, v) { 1 } else { 0 });
        }
    }
}
