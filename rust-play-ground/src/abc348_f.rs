// 1.u64 比 uszie 快
// 2.变量遮蔽比重新定义变量快很多
// 3.avx2加速
// https://atcoder.jp/contests/abc348/submissions/53822680

// 8e9 -> 300ms

use proconio::*;

fn main() {
    unsafe { run() }
}

#[target_feature(enable = "avx2")]
unsafe fn run() {
    input! {
      n: usize,
      m: usize,
      mat: [[usize; m]; n]
    }

    let mut res = 0;
    for (i, x) in mat.iter().enumerate() {
        for y in mat.iter().take(i) {
            let mut flag = false;
            for (x, y) in x.iter().zip(y.iter()) {
                flag ^= *x == *y;
            }
            res += flag as u64;
        }
    }
    println!("{}", res);
}
