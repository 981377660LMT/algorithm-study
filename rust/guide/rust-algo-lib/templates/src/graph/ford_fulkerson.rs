//! 最大流 (Ford-Fulkerson)

pub use crate::graph::flow::*;

#[derive(Clone, Debug)]
struct Edge {
    to: usize,
    rev: usize,
    cap: u64,
    is_rev: bool,
}

#[derive(Clone)]
pub struct FordFulkerson {
    size: usize,
    edges: Vec<Vec<Edge>>,
}

impl FordFulkerson {
    fn dfs(&mut self, cur: usize, sink: usize, flow: u64, check: &mut [bool]) -> u64 {
        if cur == sink {
            flow
        } else {
            check[cur] = true;
            let n = self.edges[cur].len();

            for i in 0..n {
                let e = self.edges[cur][i].clone();
                if !check[e.to] && e.cap > 0 {
                    let d = self.dfs(e.to, sink, std::cmp::min(flow, e.cap), check);
                    if d > 0 {
                        self.edges[cur][i].cap -= d;
                        self.edges[e.to][e.rev].cap += d;
                        return d;
                    }
                }
            }
            0
        }
    }
}

impl MaxFlow for FordFulkerson {
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
        let mut ret = 0;

        loop {
            let mut check = vec![false; self.size];
            let flow = self.dfs(s, t, std::u64::MAX, &mut check);
            if flow == 0 {
                return ret;
            }
            ret += flow;
        }
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
