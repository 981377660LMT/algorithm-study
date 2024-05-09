use std::{ops::Deref, rc::Rc};

fn main() {
    rc_t();
    box_ptr();
    deref();
    foo();
}

fn rc_t() {
    use ListRc::{Cons, Nil};
    let a = Rc::new(Cons(5, Rc::new(Cons(10, Rc::new(Nil)))));
    let b = Cons(3, Rc::clone(&a));
    {
        let c = Cons(4, Rc::clone(&a));
        println!("{:?}", Rc::strong_count(&a)); // 3
    }
    println!("{:?}", Rc::strong_count(&a)); // 2
}

fn deref() {
    let x = 5;
    let y = MyBox::new(x);

    assert_eq!(5, x);
    assert_eq!(5, *y);
}

fn box_ptr() {
    use List::{Cons, Nil};
    let b = Box::new(5);
    println!("b = {}", b);
    let list = Cons(1, Box::new(Cons(2, Box::new(Cons(3, Box::new(Nil))))));
    println!("{:?}", list)
}

// enum List {
//     Cons(i32, List),
//     Nil,
// }

enum ListRc {
    Cons(i32, Rc<ListRc>),
    Nil,
}

#[derive(Debug)]
enum List {
    Cons(i32, Box<List>),
    Nil,
}

struct MyBox<T>(T);

impl<T> MyBox<T> {
    fn new(x: T) -> MyBox<T> {
        MyBox(x)
    }
}

impl<T> Deref for MyBox<T> {
    type Target = T;

    fn deref(&self) -> &T {
        &self.0
    }
}

use std::cell::RefCell;

#[derive(Debug)]
struct Node {
    value: i32,
    children: RefCell<Vec<Rc<Node>>>,
}

fn foo() {
    let leaf = Rc::new(Node {
        value: 3,
        children: RefCell::new(vec![]),
    });

    let branch = Rc::new(Node {
        value: 5,
        children: RefCell::new(vec![Rc::clone(&leaf)]),
    });

    *leaf.children.borrow_mut() = vec![Rc::clone(&branch)];

    println!("{:?}", branch);
}
