//! 最大流 (Dinic)

pub use crate::graph::flow::*;
use std::{cmp::min, collections::VecDeque};

#[derive(Clone, Debug)]
struct Edge {
    to: usize,
    rev: usize,
    cap: u64,
    is_rev: bool,
}

#[derive(Clone)]
pub struct Dinic {
    size: usize,
    edges: Vec<Vec<Edge>>,
}

impl Dinic {
    fn dfs(&mut self, path: &mut Vec<(usize, usize)>, cur: usize, t: usize, level: &[u32]) -> u64 {
        if cur == t {
            let mut f = std::u64::MAX;

            for (i, j) in path.clone() {
                let e = self.edges[i][j].clone();
                f = min(f, e.cap);
            }

            for &mut (i, j) in path {
                let e = &mut self.edges[i][j];
                e.cap -= f;
                let Edge { to, rev, .. } = *e;
                self.edges[to][rev].cap += f;
            }

            f
        } else {
            let n = self.edges[cur].len();
            let mut f = 0;

            for i in 0..n {
                let e = self.edges[cur][i].clone();
                if e.cap > 0 && level[e.to] > level[cur] {
                    path.push((cur, i));
                    f += self.dfs(path, e.to, t, level);
                    path.pop();
                }
            }
            f
        }
    }
}

impl MaxFlow for Dinic {
    type Cap = u64;
    fn new(size: usize) -> Self {
        Self {
            size,
            edges: vec![vec![]; size],
        }
    }

    fn add_edge(&mut self, u: usize, v: usize, cap: Self::Cap) {
        let rev = self.edges[v].len();
        self.edges[u].push(Edge {
            to: v,
            rev,
            cap,
            is_rev: false,
        });
        let rev = self.edges[u].len() - 1;
        self.edges[v].push(Edge {
            to: u,
            rev,
            cap: 0,
            is_rev: true,
        });
    }

    fn max_flow(&mut self, s: usize, t: usize) -> Self::Cap {
        let mut f = 0;
        loop {
            let mut level = vec![0; self.size];
            level[s] = 1;
            let mut q = VecDeque::new();
            q.push_back(s);

            while let Some(cur) = q.pop_front() {
                for e in &self.edges[cur] {
                    if level[e.to] == 0 && e.cap > 0 {
                        level[e.to] = level[cur] + 1;
                        q.push_back(e.to);
                    }
                }
            }

            if level[t] == 0 {
                break;
            }

            f += self.dfs(&mut vec![], s, t, &level);
        }
        f
    }

    fn get_edges(&self, i: usize) -> Vec<(usize, u64)> {
        self.edges[i]
            .iter()
            .filter(|e| !e.is_rev)
            .map(|e| (e.to, e.cap))
            .collect()
    }

    fn reset(&mut self) {
        todo!();
    }
}
