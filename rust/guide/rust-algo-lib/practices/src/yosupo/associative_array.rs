use std::collections::HashMap;

use templates::misc::io::run_with_io;

#[allow(unused)]
fn main() {
    run_with_io(|input, ouput| {
        let q: usize = input.next();
        let mut map = HashMap::new();
        for _ in 0..q {
            let op: u8 = input.next();
            if op == 0 {
                let k: i64 = input.next();
                let v: i64 = input.next();
                map.insert(k, v);
            } else {
                let k: i64 = input.next();
                let res = map.get(&k).unwrap_or(&0);
                writeln!(ouput, "{}", res).ok();
            }
        }
    })
}
