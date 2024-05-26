pub trait MaxFlow {
    type Cap;
    fn new(n: usize) -> Self;
    fn add_edge(&mut self, u: usize, v: usize, cap: Self::Cap);
    fn max_flow(&mut self, s: usize, t: usize) -> Self::Cap;
    fn get_edges(&self, i: usize) -> Vec<(usize, Self::Cap)>;
    fn reset(&mut self);
}
