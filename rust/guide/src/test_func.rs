fn main() {}

fn max<T: PartialOrd + Copy>(list: &[T]) -> T {
    let mut max: T = &list[0];
    for &item in list.iter() {
        if item > max {
            max = item;
        }
    }
    max
}

struct Point<T> {
    x: T,
    y: T,
}
