pub struct UnionFind {
    part: usize,
    parent_or_size: Vec<i32>,
}

impl UnionFind {
    pub fn new(n: usize) -> Self {
        Self {
            part: n,
            parent_or_size: vec![-1; n],
        }
    }

    /// # Arguments
    ///
    /// - `pre_merge` - The callback function before merging.
    ///
    /// # Returns
    ///
    /// - `bool` - Whether union successfully.
    ///
    /// # Example
    ///
    /// ```
    /// use templates::ds::union_find::UnionFind;
    /// let mut uf = UnionFind::new(10);
    /// uf.union(0, 1, None);
    /// uf.union(0, 2, Some(|big, small| println!("{} {}", big, small)));
    /// ```
    pub fn union(
        &mut self,
        a: usize,
        b: usize,
        pre_merge: Option<fn(big: usize, small: usize)>,
    ) -> bool {
        let (mut x, mut y) = (self.find(a), self.find(b));
        if x == y {
            return false;
        }
        if -self.parent_or_size[x] < -self.parent_or_size[y] {
            std::mem::swap(&mut x, &mut y);
        }
        if let Some(f) = pre_merge {
            f(x, y);
        }
        self.part -= 1;
        self.parent_or_size[x] += self.parent_or_size[y];
        self.parent_or_size[y] = x as i32;
        true
    }

    pub fn find(&mut self, a: usize) -> usize {
        if self.parent_or_size[a] < 0 {
            return a;
        }
        self.parent_or_size[a] = self.find(self.parent_or_size[a] as usize) as i32;
        self.parent_or_size[a] as usize
    }

    pub fn is_connected(&mut self, a: usize, b: usize) -> bool {
        self.find(a) == self.find(b)
    }

    pub fn size(&mut self, a: usize) -> usize {
        let x = self.find(a);
        -self.parent_or_size[x] as usize
    }

    pub fn groups(&mut self) -> Vec<Vec<usize>> {
        let n = self.parent_or_size.len();
        let mut groups = vec![vec![]; n];
        for x in 0..n {
            groups[self.find(x)].push(x);
        }
        groups.into_iter().filter(|x| !x.is_empty()).collect()
    }

    pub fn part(&self) -> usize {
        self.part
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_union_find() {
        let mut uf = UnionFind::new(4);
        uf.union(0, 1, None);
        assert!(uf.is_connected(0, 1));
        assert_eq!(uf.size(0), 2);
        uf.union(1, 2, None);
        assert!(uf.is_connected(0, 2));
        assert_eq!(uf.size(0), 3);
        assert!(!uf.is_connected(0, 3));
        assert_eq!(uf.size(3), 1);
        assert_eq!(uf.groups(), vec![vec![0, 1, 2], vec![3]]);
    }
}
