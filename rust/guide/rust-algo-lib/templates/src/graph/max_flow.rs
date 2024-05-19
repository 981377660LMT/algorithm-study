// ---------- begin max flow (Dinic) ----------
mod maxflow {
    pub trait MaxFlowCapacity:
        Copy + Ord + std::ops::Add<Output = Self> + std::ops::Sub<Output = Self>
    {
        fn zero() -> Self;
        fn inf() -> Self;
    }

    macro_rules! impl_primitive_integer_capacity {
        ($x:ty, $y:expr) => {
            impl MaxFlowCapacity for $x {
                fn zero() -> Self {
                    0
                }
                fn inf() -> Self {
                    $y
                }
            }
        };
    }

    impl_primitive_integer_capacity!(u32, std::u32::MAX);
    impl_primitive_integer_capacity!(u64, std::u64::MAX);
    impl_primitive_integer_capacity!(i32, std::i32::MAX);
    impl_primitive_integer_capacity!(i64, std::i64::MAX);

    #[derive(Clone)]
    struct Edge<Cap> {
        to_: u32,
        inv_: u32,
        cap_: Cap,
    }

    impl<Cap> Edge<Cap> {
        fn new(to: usize, inv: usize, cap: Cap) -> Self {
            Edge {
                to_: to as u32,
                inv_: inv as u32,
                cap_: cap,
            }
        }
        fn to(&self) -> usize {
            self.to_ as usize
        }
        fn inv(&self) -> usize {
            self.inv_ as usize
        }
    }

    impl<Cap: MaxFlowCapacity> Edge<Cap> {
        fn add(&mut self, cap: Cap) {
            self.cap_ = self.cap_ + cap;
        }
        fn sub(&mut self, cap: Cap) {
            self.cap_ = self.cap_ - cap;
        }
        fn cap(&self) -> Cap {
            self.cap_
        }
    }

    pub struct Graph<Cap> {
        graph: Vec<Vec<Edge<Cap>>>,
    }

    #[allow(dead_code)]
    pub struct EdgeIndex {
        src: usize,
        dst: usize,
        x: usize,
        y: usize,
    }

    impl<Cap: MaxFlowCapacity> Graph<Cap> {
        pub fn new(size: usize) -> Self {
            Self {
                graph: vec![vec![]; size],
            }
        }
        pub fn add_edge(&mut self, src: usize, dst: usize, cap: Cap) -> EdgeIndex {
            assert!(src.max(dst) < self.graph.len());
            assert!(cap >= Cap::zero());
            assert!(src != dst);
            let x = self.graph[src].len();
            let y = self.graph[dst].len();
            self.graph[src].push(Edge::new(dst, y, cap));
            self.graph[dst].push(Edge::new(src, x, Cap::zero()));
            EdgeIndex { src, dst, x, y }
        }
        // src, dst, used, intial_capacity
        #[allow(dead_code)]
        pub fn get_edge(&self, e: &EdgeIndex) -> (usize, usize, Cap, Cap) {
            let max = self.graph[e.src][e.x].cap() + self.graph[e.dst][e.y].cap();
            let used = self.graph[e.dst][e.y].cap();
            (e.src, e.dst, used, max)
        }
        pub fn flow(&mut self, src: usize, dst: usize) -> Cap {
            let size = self.graph.len();
            assert!(src.max(dst) < size);
            assert!(src != dst);
            let mut queue = std::collections::VecDeque::new();
            let mut level = vec![0; size];
            let mut it = vec![0; size];
            let mut ans = Cap::zero();
            loop {
                (|| {
                    level.clear();
                    level.resize(size, 0);
                    level[src] = 1;
                    queue.clear();
                    queue.push_back(src);
                    while let Some(v) = queue.pop_front() {
                        let d = level[v] + 1;
                        for e in self.graph[v].iter() {
                            let u = e.to();
                            if e.cap() > Cap::zero() && level[u] == 0 {
                                level[u] = d;
                                if u == dst {
                                    return;
                                }
                                queue.push_back(u);
                            }
                        }
                    }
                })();
                if level[dst] == 0 {
                    break;
                }
                it.clear();
                it.resize(size, 0);
                loop {
                    let f = self.dfs(dst, src, Cap::inf(), &mut it, &level);
                    if f == Cap::zero() {
                        break;
                    }
                    ans = ans + f;
                }
            }
            ans
        }
        fn dfs(&mut self, v: usize, src: usize, cap: Cap, it: &mut [usize], level: &[u32]) -> Cap {
            if v == src {
                return cap;
            }
            while let Some((u, inv)) = self.graph[v].get(it[v]).map(|p| (p.to(), p.inv())) {
                if level[u] + 1 == level[v] && self.graph[u][inv].cap() > Cap::zero() {
                    let cap = cap.min(self.graph[u][inv].cap());
                    let c = self.dfs(u, src, cap, it, level);
                    if c > Cap::zero() {
                        self.graph[v][it[v]].add(c);
                        self.graph[u][inv].sub(c);
                        return c;
                    }
                }
                it[v] += 1;
            }
            Cap::zero()
        }
    }
}
// ---------- end max flow (Dinic) ----------
