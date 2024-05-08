fn main() {
    let v1 = vec![1, 2, 3, 4, 5];
    let mut v1_iter = v1.iter();

    assert_eq!(v1_iter.next(), Some(&1));
    assert_eq!(v1_iter.next(), Some(&2));
    assert_eq!(v1_iter.next(), Some(&3));

    iterator_sum();
    map();
}

fn iterator_sum() {
    let nums = vec![1, 2, 3, 4, 5];
    let sum: i32 = nums.iter().sum();
    assert_eq!(sum, 15);
}

fn map() {
    let nums = vec![1, 2, 3, 4, 5];
    let plus_one: Vec<i32> = nums.iter().map(|x| x + 1).collect();
    assert_eq!(plus_one, vec![2, 3, 4, 5, 6]);
}

struct Counter {
    count: u32,
}

impl Counter {
    fn new() -> Counter {
        Counter { count: 0 }
    }
}

impl Iterator for Counter {
    type Item = u32;

    fn next(&mut self) -> Option<Self::Item> {
        if self.count < 5 {
            self.count += 1;
            Some(self.count)
        } else {
            None
        }
    }
}
