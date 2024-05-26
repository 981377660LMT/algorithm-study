//! Project Selection Problem
//! 燃やす埋める問題

use crate::graph::flow::*;

/// Project Selection Problem
///
/// # Problems
/// - [ARC 085 E - MUL](https://atcoder.jp/contests/arc085/tasks/arc085_c)
/// - [AOJ 3058 Ghost](http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=3058)
/// - [AOJ 2903 Board](http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2903)
/// - [ABC 193 F - Zebraness](https://atcoder.jp/contests/abc193/tasks/abc193_f)
///
/// # References
/// - [https://ferin-tech.hatenablog.com/entry/2019/10/28/燃やす埋める問題](https://ferin-tech.hatenablog.com/entry/2019/10/28/%E7%87%83%E3%82%84%E3%81%99%E5%9F%8B%E3%82%81%E3%82%8B%E5%95%8F%E9%A1%8C)
/// - [https://kmyk.github.io/blog/blog/2017/12/05/minimum-cut-and-project-selection-problem/](https://kmyk.github.io/blog/blog/2017/12/05/minimum-cut-and-project-selection-problem/)
///
/// # Verification
/// | function | verify |
/// | -------- | ------ |
/// | `penalty_if_red` | |
/// | `penalty_if_blue` | [ARC085 #27489484](https://atcoder.jp/contests/arc085/submissions/27489484) |
/// | `gain_if_red` | |
/// | `gain_if_blue` | [ARC085 #27489484](https://atcoder.jp/contests/arc085/submissions/27489484) |
/// | `penalty_if_red_blue` | |
/// | `penalty_if_different` | |
/// | `must_be_red` | |
/// | `must_be_blue` | |
/// | `if_red_then_must_be_red` | [ARC085 #27489484](https://atcoder.jp/contests/arc085/submissions/27489484) |
/// | `gain_if_both_red` | |
/// | `gain_if_both_blue` | |

#[derive(Clone)]
pub struct PSP {
    size: usize,
    src: usize,
    sink: usize,
    edges: Vec<(usize, usize, u64)>,
    node_count: usize,
    default_gain: u64,
}

impl PSP {
    pub fn new(size: usize) -> Self {
        Self {
            size,
            src: size,
            sink: size + 1,
            edges: vec![],
            node_count: size + 2,
            default_gain: 0,
        }
    }

    /// 頂点iが<font color="red"><b>赤</b></font>ならばcの損失になる。
    pub fn penalty_if_red(&mut self, i: usize, c: u64) {
        assert!(i < self.size);
        self.edges.push((i, self.sink, c));
    }

    /// 頂点iが<font color="blue"><b>青</b></font>ならばcの損失になる。
    pub fn penalty_if_blue(&mut self, i: usize, c: u64) {
        assert!(i < self.size);
        self.edges.push((self.src, i, c));
    }

    /// 頂点iが<font color="red"><b>赤</b></font>ならばcの利益を得る。
    pub fn gain_if_red(&mut self, i: usize, c: u64) {
        assert!(i < self.size);
        self.default_gain += c;
        self.penalty_if_blue(i, c);
    }

    /// 頂点iが<font color="blue"><b>青</b></font>ならばcの利益を得る。
    pub fn gain_if_blue(&mut self, i: usize, c: u64) {
        assert!(i < self.size);
        self.default_gain += c;
        self.penalty_if_red(i, c);
    }

    /// 頂点iが<font color="red"><b>赤</b></font>かつ頂点jが<font color="blue"><b>青</b></font>ならばcの損失となる。
    pub fn penalty_if_red_blue(&mut self, i: usize, j: usize, c: u64) {
        assert!(i < self.size && j < self.size);
        self.edges.push((i, j, c));
    }

    /// 頂点iとjが異なる色ならばcの損失となる。
    pub fn penalty_if_different(&mut self, i: usize, j: usize, c: u64) {
        assert!(i < self.size && j < self.size);
        self.edges.push((i, j, c));
        self.edges.push((j, i, c));
    }

    /// 頂点iは<font color="red"><b>赤</b></font>でなければならない。
    pub fn must_be_red(&mut self, i: usize) {
        assert!(i < self.size);
        self.penalty_if_blue(i, std::u64::MAX);
    }

    /// 頂点iは<font color="blue"><b>青</b></font>でなければならない。
    pub fn must_be_blue(&mut self, i: usize) {
        assert!(i < self.size);
        self.penalty_if_red(i, std::u64::MAX);
    }

    /// 頂点iが<font color="red"><b>赤</b></font>ならば、頂点jも<font color="red"><b>赤</b></font>でなければならない。
    pub fn if_red_then_must_be_red(&mut self, i: usize, j: usize) {
        assert!(i < self.size && j < self.size);
        self.penalty_if_red_blue(i, j, std::u64::MAX);
    }

    /// 頂点iとjがともに<font color="red"><b>赤</b></font>ならばcの利益を得る。
    pub fn gain_if_both_red(&mut self, i: usize, j: usize, c: u64) {
        assert!(i < self.size && j < self.size);
        self.default_gain += c;
        let w = self.node_count;
        self.node_count += 1;

        self.edges.push((self.src, w, c));
        self.edges.push((w, i, std::u64::MAX));
        self.edges.push((w, j, std::u64::MAX));
    }

    /// 頂点iとjがともに<font color="blue"><b>青</b></font>ならばcの利益を得る。
    pub fn gain_if_both_blue(&mut self, i: usize, j: usize, c: u64) {
        assert!(i < self.size && j < self.size);
        self.default_gain += c;
        let w = self.node_count;
        self.node_count += 1;

        self.edges.push((w, self.sink, c));
        self.edges.push((i, w, std::u64::MAX));
        self.edges.push((j, w, std::u64::MAX));
    }

    /// must be制約を破った場合、`None`を返す。そうでなければ、利益の最大値を`Some`に包んで返す。
    pub fn solve<F: MaxFlow<Cap = u64>>(self) -> Option<i64> {
        let mut f = F::new(self.node_count);
        for (u, v, c) in self.edges {
            f.add_edge(u, v, c);
        }

        let flow = f.max_flow(self.src, self.sink);

        if flow == std::u64::MAX {
            None
        } else {
            Some(self.default_gain as i64 - flow as i64)
        }
    }
}
