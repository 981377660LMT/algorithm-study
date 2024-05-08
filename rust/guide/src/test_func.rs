fn main() {
    let p: Point<i32> = Point { x: 1, y: 2 };
    p.foo();
    p.bar();
    p.summarize();

    let list = vec![1, 2, 3, 4, 5];
    println!("max = {}", max(&list));
}

fn max<T: PartialOrd + Clone>(list: &[T]) -> &T {
    let mut max: &T = &list[0];
    for v in list.iter() {
        if v > max {
            max = v;
        }
    }
    max
}

struct Point<T> {
    x: T,
    y: T,
}

impl<T> Point<T> {
    fn foo(&self) {}
}

impl Point<i32> {
    fn bar(&self) {}
}

pub trait Summary {
    fn summarize(&self) -> String;
}

impl Summary for Point<i32> {
    fn summarize(&self) -> String {
        format!("Point: x={}, y={}", self.x, self.y)
    }
}
