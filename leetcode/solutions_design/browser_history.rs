// struct BrowserHistory {
//     fifo: Vec<String>,
//     cur: usize,
// }

// /**
//  * `&self` means the method takes an immutable reference.
//  * If you need a mutable reference, change it to `&mut self` instead.
//  */
// impl BrowserHistory {

//     fn new(homepage: String) -> Self {
//         let mut fifo = Vec::with_capacity(2048);
//         fifo.push(homepage);
//         Self {fifo, cur: 0}
//     }

//     fn visit(&mut self, url: String) {
//         self.fifo.truncate(self.cur + 1);
//         self.fifo.push(url);
//         self.cur += 1;
//     }

//     fn back(&mut self, steps: i32) -> String {
//         let steps = steps as usize;
//         self.cur = if steps > self.cur {0}
//         else {self.cur - steps};
//         self.fifo[self.cur].clone()
//     }

//     fn forward(&mut self, steps: i32) -> String {
//         let steps = steps as usize;
//         self.cur = (steps + self.cur).min(self.fifo.len() - 1);
//         self.fifo[self.cur].clone()
//     }
// }

// 作者：934786601
// 链接：https://leetcode.cn/problems/design-browser-history/solutions/
// 来源：力扣（LeetCode）
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

struct BrowserHistory {
    ptr: usize,
    history: Vec<String>,
}

impl BrowserHistory {
    fn new(homepage: String) -> Self {
        let mut history = Vec::with_capacity(2048);
        history.push(homepage);
        Self { ptr: 0, history }
    }

    fn visit(&mut self, url: String) {
        self.history.truncate(self.ptr + 1);
        self.history.push(url);
        self.ptr += 1;
    }

    fn back(&mut self, steps: i32) -> String {
        let steps = steps as usize;
        self.ptr = if steps > self.ptr {
            0
        } else {
            self.ptr - steps
        };
        self.history[self.ptr].to_owned()
    }

    fn forward(&mut self, steps: i32) -> String {
        let steps = steps as usize;
        self.ptr = (steps + self.ptr).min(self.history.len() - 1);
        self.history[self.ptr].to_owned()
    }
}
