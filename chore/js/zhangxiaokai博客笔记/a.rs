#[derive(Debug, Copy, Clone)]
struct Meter {
    value: f64,
}

impl Meter {
    fn new(value: f64) -> Self {
        Self { value }
    }
}
