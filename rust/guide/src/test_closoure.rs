use std::thread;
use std::time::Duration;

fn main() {
    generate_workout(3, 45);
    closure_with_move();
}

fn generate_workout(intensity: u32, random_number: u32) {
    // useCallback
    let mut expensive_closure = Lazy::new(|num| {
        println!("calculating slowly...");
        thread::sleep(Duration::from_secs(2));
        num
    });

    if intensity < 25 {
        println!("Today, do {} pushups!", expensive_closure.value(intensity));
        println!("Next, do {} situps!", expensive_closure.value(intensity));
    } else {
        if random_number == 3 {
            println!("Take a break today! Remember to stay hydrated!");
        } else {
            println!(
                "Today, run for {} minutes!",
                expensive_closure.value(intensity)
            );
        }
    }
}

struct Lazy<T: Fn(u32) -> u32> {
    calculation: T,
    value: Option<u32>,
}

impl<T: Fn(u32) -> u32> Lazy<T> {
    fn new(calculation: T) -> Lazy<T> {
        Lazy {
            calculation,
            value: None,
        }
    }

    fn value(&mut self, arg: u32) -> u32 {
        match self.value {
            Some(v) => v,
            None => {
                let v = (self.calculation)(arg);
                self.value = Some(v);
                v
            }
        }
    }
}

fn closure_with_move() {
    let x = vec![1, 2, 3];
    let equal_to_x = move |z| z == x;
    // println!("can't use x here: {:?}", x);
    let y = vec![1, 2, 3];
    assert!(equal_to_x(y));
}
