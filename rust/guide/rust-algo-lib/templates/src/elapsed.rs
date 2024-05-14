#[macro_export]
macro_rules! elapsed {
    ($e: expr) => {
        use std::time::Instant;
        let now = Instant::now();
        {
            $e;
        }
        println!("elapsed: {:?}", now.elapsed());
    };
}
