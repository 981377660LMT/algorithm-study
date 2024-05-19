fn main() {
    println!("Hello, union_find!sdvfc bv")
}

pub struct UnionFind {
    parent: Vec<i32>,
    part: usize,
}

impl UnionFind {
    pub fn new(n: usize) -> Self {
        let parent = vec![0; n];

        Self { parent, part: n }
    }
}
