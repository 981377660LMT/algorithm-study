struct MyCircularQueue {
    rear: usize,
    size: usize,
    arr: Vec<i32>,
}

impl MyCircularQueue {
    fn new(k: i32) -> Self {
        MyCircularQueue {
            rear: 0,
            size: 0,
            arr: vec![0; k as usize],
        }
    }

    fn en_queue(&mut self, value: i32) -> bool {
        if self.is_full() {
            return false;
        }
        self.arr[self.rear] = value;
        self.size += 1;
        self.rear += 1;
        if self.rear == self.arr.len() {
            self.rear = 0;
        }
        true
    }

    fn de_queue(&mut self) -> bool {
        if self.is_empty() {
            return false;
        }
        self.size -= 1;
        true
    }

    fn front(&self) -> i32 {
        if self.is_empty() {
            return -1;
        }
        self.arr[(self.rear - self.size + self.arr.len()) % self.arr.len()]
    }

    fn rear(&self) -> i32 {
        if self.is_empty() {
            return -1;
        }
        self.arr[(self.rear - 1 + self.arr.len()) % self.arr.len()]
    }

    fn is_empty(&self) -> bool {
        self.size == 0
    }

    fn is_full(&self) -> bool {
        self.size == self.arr.len()
    }
}
