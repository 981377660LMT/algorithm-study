// !在Rust中，as_ref().unwrap().borrow()这一连串调用通常出现在处理Option<Rc<RefCell<T>>>类型时.
// 尤其是在需要内部可变性和引用计数共享所有权的场景中.
// rust实现链表.

use std::cell::RefCell;
use std::rc::Rc;

struct ListNode {
    val: i32,
    prev: Option<Rc<RefCell<ListNode>>>,
    next: Option<Rc<RefCell<ListNode>>>,
}

impl ListNode {
    fn new(
        val: i32,
        prev: Option<Rc<RefCell<ListNode>>>,
        next: Option<Rc<RefCell<ListNode>>>,
    ) -> ListNode {
        ListNode { prev, next, val }
    }
}

struct MyLinkedList {
    size: i32,
    head: Option<Rc<RefCell<ListNode>>>,
}

impl MyLinkedList {
    fn new() -> Self {
        let head = Rc::new(RefCell::new(ListNode::new(-1, None, None)));
        head.borrow_mut().prev = Some(head.clone());
        head.borrow_mut().next = Some(head.clone());
        MyLinkedList {
            size: 0,
            head: Some(head),
        }
    }

    fn get(&self, index: i32) -> i32 {
        if index < 0 || index >= self.size {
            return -1;
        }
        self.find_node(index).as_ref().unwrap().borrow().val
    }

    fn add_at_head(&mut self, val: i32) {
        let next = self.head.as_ref().unwrap().borrow().next.clone();
        let cur = Rc::new(RefCell::new(ListNode::new(
            val,
            self.head.clone(),
            next.clone(),
        )));
        next.as_ref().unwrap().borrow_mut().prev = Some(cur.clone());
        self.head.as_ref().unwrap().borrow_mut().next = Some(cur.clone());
        self.size += 1;
    }

    fn add_at_tail(&mut self, val: i32) {
        let prev = self.head.as_ref().unwrap().borrow().prev.clone();
        let cur = Rc::new(RefCell::new(ListNode::new(
            val,
            prev.clone(),
            self.head.clone(),
        )));
        prev.as_ref().unwrap().borrow_mut().next = Some(cur.clone());
        self.head.as_ref().unwrap().borrow_mut().prev = Some(cur);
        self.size += 1;
    }

    fn add_at_index(&mut self, index: i32, val: i32) {
        if index > self.size {
            return;
        }
        if index < 0 {
            self.add_at_head(val);
            return;
        }
        if index == self.size {
            self.add_at_tail(val);
            return;
        }
        let cur = self.find_node(index);
        let new_node = Rc::new(RefCell::new(ListNode::new(
            val,
            cur.as_ref().unwrap().borrow().prev.clone(),
            cur.clone(),
        )));
        cur.as_ref()
            .unwrap()
            .borrow_mut()
            .prev
            .as_ref()
            .unwrap()
            .borrow_mut()
            .next = Some(new_node.clone());
        cur.as_ref().unwrap().borrow_mut().prev = Some(new_node);
        self.size += 1;
    }

    fn delete_at_index(&mut self, index: i32) {
        if self.size <= 0 || index < 0 || index >= self.size {
            return;
        }
        if self.size == 1 {
            self.head.as_ref().unwrap().borrow_mut().prev = self.head.clone();
            self.head.as_ref().unwrap().borrow_mut().next = self.head.clone();
            self.size -= 1;
            return;
        }
        let cur = self.find_node(index);
        let (prev, next) = (
            cur.as_ref().unwrap().borrow().prev.clone(),
            cur.as_ref().unwrap().borrow().next.clone(),
        );
        prev.as_ref().unwrap().borrow_mut().next = next.clone();
        next.as_ref().unwrap().borrow_mut().prev = prev.clone();
        self.size -= 1;
    }

    fn find_node(&self, index: i32) -> Option<Rc<RefCell<ListNode>>> {
        if self.size == 0 {
            return self.head.as_ref().unwrap().borrow().next.clone();
        }
        if self.size == index {
            return self.head.as_ref().unwrap().borrow().prev.clone();
        }
        let mut cur = self.head.as_ref().unwrap().borrow().next.clone();
        for _ in 0..index {
            let tmp = cur.as_ref().unwrap().borrow().next.clone();
            cur = tmp;
        }
        cur
    }
}
