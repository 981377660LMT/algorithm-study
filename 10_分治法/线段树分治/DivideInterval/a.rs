#[allow(unused_imports)]
use ac_library::*;
#[allow(unused_imports)]
use proconio::{marker::*, *};
#[allow(unused_imports)]
use std::collections::*;

type Mint = ModInt;

fn main() {
    ModInt::set_modulus(100);

    input_interactive! {
        _: usize,
        mut l: usize,
        r: usize,
    }
    let mut r = r + 1;
    let mut i = 0;
    let mut ans = Mint::raw(0);
    while l != r {
        if l & 1 == 1 {
            if l & 2 == 0 && l + 1 != r && l + 2 != r {
                l -= 1;
                println!("? {i} {l}");
                input_interactive! { q: Mint }
                ans -= q;
            } else {
                println!("? {i} {l}");
                input_interactive! { q: Mint }
                ans += q;
                l += 1;
            }
        }
        if r & 1 == 1 {
            if r & 2 == 0 || l + 1 == r {
                r -= 1;
                println!("? {i} {r}");
                input_interactive! { q: Mint }
                ans += q;
            } else {
                println!("? {i} {r}");
                input_interactive! { q: Mint }
                ans -= q;
                r += 1;
            }
        }
        l >>= 1;
        r >>= 1;
        i += 1;
    }
    println!("! {ans}");
}
