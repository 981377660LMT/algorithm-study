use std::{sync::mpsc, thread, time::Duration};

fn main() {
    let v = vec![1, 2, 3];
    let join_handle = thread::spawn(move || {
        for i in 1..10 {
            println!("hi number {} from the spawned thread!", i);
            thread::sleep(Duration::from_millis(1));
        }
        println!("{:?}", v)
    });
    // drop(v); // v is moved to the spawned thread, so it's not available here anymore

    for i in 1..5 {
        println!("hi number {} from the main thread!", i);
        thread::sleep(Duration::from_millis(1));
    }

    join_handle.join().unwrap();

    test_channel();
}

fn test_channel() {
    let (sender, receiver) = mpsc::channel();
    thread::spawn(move || {
        let val = String::from("hi");
        sender.send(val).unwrap();
        // println!("val is {}", val);
    });

    for v in receiver {
        println!("Got: {}", v);
    }
}
