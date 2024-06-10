// 为什么要用裸指针：
// 在双向链表中，使用 Box 是行不通的，因为会出现一个节点同时存在多个可变引用的情况；
// 所以需要使用 unsafe 来手动管理，所以需要使用裸指针
//
// 裸指针 *mut T 和 NonNull 有什么区别呢？
// !NonNull 提供了比 *mut T 更多的内容：支持协变类型、空指针优化，并且可以保证指针非空；

use std::{error::Error, fmt::Debug, marker::PhantomData, mem, ptr::NonNull};

struct Node<T> {
    // 如果存在这个 Node，则该Node中必定是有值的，
    // 保证不会出现 Node 存在而 val 为 null 的情况
    val: T,
    next: Option<NonNull<Node<T>>>,
    prev: Option<NonNull<Node<T>>>,
}

impl<T> Node<T> {
    fn new(val: T) -> Node<T> {
        Node {
            val,
            prev: None,
            next: None,
        }
    }

    fn into_val(self: Box<Self>) -> T {
        self.val
    }
}

pub struct LinkedList<T> {
    length: usize,
    head: Option<NonNull<Node<T>>>,
    tail: Option<NonNull<Node<T>>>,
}

impl<T> LinkedList<T> {
    fn new() -> Self {
        Self {
            length: 0,
            head: None,
            tail: None,
        }
    }

    pub fn push_front(&mut self, val: T) {
        // 使用 Box 包装（这么做是方便我们后面直接从 Box 获取到裸指针
        let mut node = Box::new(Node::new(val));
        node.next = self.head;
        node.prev = None;

        // 使用 Box::into_raw(node)，将 node 转为裸指针
        // 当对某个被 Box 包装的变量调用了 Box::into_raw 后，
        // 变量将会被转化为裸指针，同时指针指向的内存的管理权会被交给我们自己
        // !Box::into_raw 所做的其实就是消费掉 Box 并返回指针，并且保证不会像 Box 退出作用域后释放指针指向的内存
        // 因此，需要我们自己保存这个裸指针，并在适当时候释放这个裸指针指向的内存
        // 那么如何释放由 Box 转换所得的裸指针指向的内存呢？
        // !最简单的方法是使用 Box::from_raw 函数将原始指针转换回 Box，从而允许 Box 析构函数执行清理；
        // 所以我们只需要将裸指针再转为实际的 Box，然后通过 Box 退出作用域后直接释放内存即可；
        let node = NonNull::new(Box::into_raw(node));

        match self.head {
            None => self.tail = node,
            Some(head) => unsafe {
                // 将链表中的 head 裸指针进行解引用并修改其 prev 值
                // 解引用裸指针必须unsafe
                (*head.as_ptr()).prev = node;
            },
        }

        self.head = node;
        self.length += 1;
    }

    pub fn push_back(&mut self, val: T) {
        // use box to help generate raw ptr
        let mut node = Box::new(Node::new(val));
        node.next = None;
        node.prev = self.tail;

        let node = NonNull::new(Box::into_raw(node));
        match self.tail {
            None => self.head = node,
            // Not creating new mutable (unique!) references overlapping `element`.
            Some(tail) => unsafe {
                (*tail.as_ptr()).next = node;
            },
        }

        self.tail = node;
        self.length += 1;
    }

    // 所谓弹出就是：将元素从链表中删除，并且返回具有所有权的 T（如果存在的话）；
    // ps:在 Rust 中所有的变量一定都不为 Null，即不会发生空指针；
    // Null 值的语义就是通过枚举类型 Option 来显示标注的！
    pub fn pop_front(&mut self) -> Option<T> {
        // 注意到代码风格，只是调用了 self.head.map() 即完成了所有功能；
        // !option.map函数可以忽略掉 None 的情况，只处理 Some 的情况；
        self.head.map(|node| {
            self.length -= 1;

            unsafe {
                // 使用 Box::from_raw 将裸指针还原为 Box<Node<T>> 类型（为返回头节点数据做准备）；
                let node = Box::from_raw(node.as_ptr());
                self.head = node.next;
                match self.head {
                    None => self.tail = None,
                    Some(next) => {
                        (*next.as_ptr()).prev = None;
                    }
                }

                node.into_val()
            }
        })
    }

    pub fn pop_back(&mut self) -> Option<T> {
        self.tail.map(|node| {
            self.length -= 1;

            unsafe {
                let node = Box::from_raw(node.as_ptr());
                self.tail = node.prev;
                match self.tail {
                    None => self.head = None,
                    Some(prev) => {
                        (*prev.as_ptr()).next = None;
                    }
                }

                node.into_val()
            }
        })
    }

    // 仅仅返回元素的引用，而元素的所有权还是在链表中
    pub fn peek_front(&self) -> Option<&T> {
        // Option 提供了 as_ref 方法，可以将 Option<T> 转为 Option<&T> 而不用频繁的拆包再包装；
        // 这里直接支持这个操作的原因是因为：
        // 我们使用了 NonNull 类型，保证了指针一定不为空，即：裸指针一定不为空指针！
        unsafe { self.head.as_ref().map(|node| &(node.as_ref().val)) }
    }

    pub fn peek_back(&self) -> Option<&T> {
        unsafe { self.tail.as_ref().map(|node| &(node.as_ref().val)) }
    }

    pub fn peek_front_mut(&mut self) -> Option<&mut T> {
        // 在 Rust 中，如果修改一个容器中的元素，首先这个容器需要是可变的！
        unsafe { self.head.as_mut().map(|node| &mut (node.as_mut().val)) }
    }

    pub fn peek_back_mut(&mut self) -> Option<&mut T> {
        unsafe { self.tail.as_mut().map(|node| &mut (node.as_mut().val)) }
    }

    // 如果仅仅返回 None，api调用方不能确定是因为 index 传错而导致的 None，还是链表本身就是空的
    pub fn get_by_index(&self, index: usize) -> Result<Option<&T>, Box<dyn Error>> {
        let len = self.length;
        if index >= len {
            return Err(Box::new(IndexOutOfRangeError));
        }

        let offset_from_end = len - index - 1;
        let mut cur;
        if index <= offset_from_end {
            // 从头开始遍历
            cur = self.head;
            for _ in 0..index {
                match cur.take() {
                    // Rust 中的 = 是 move 语义，这样做原链表中的 head 不就变成空值了！
                    // 注：在 Rust 中默认是 Move 语义，但是如果实现了 Copy Trait就会变为 Copy 语义；
                    // 因此，明确一个变量是否实现了 Copy Trait 是非常重要的
                    // 与Clone不同，Copy方式是隐式作用于类型变量，通过赋值语句来完成；
                    // 最终会将变量 cur 也赋值为指向链表头部的裸指针
                    None => cur = self.head,
                    Some(node) => unsafe {
                        cur = (*node.as_ptr()).next;
                    },
                }
            }
        } else {
            // 从尾开始遍历
            cur = self.tail;
            for _ in 0..offset_from_end {
                match cur.take() {
                    None => cur = self.tail,
                    Some(node) => unsafe {
                        cur = (*node.as_ptr()).prev;
                    },
                }
            }
        }

        Ok(cur.map(|node| unsafe { &(node.as_ref()).val }))
    }

    pub fn get_by_index_mut(&self, index: usize) -> Result<Option<&mut T>, Box<dyn Error>> {
        let mut cur = self._get_by_index_mut(index)?;
        Ok(cur.as_mut().map(|node| unsafe { &mut (node.as_mut().val) }))
    }

    fn _get_by_index_mut(&self, index: usize) -> Result<Option<NonNull<Node<T>>>, Box<dyn Error>> {
        let len = self.length;
        if index >= len {
            return Err(Box::new(IndexOutOfRangeError));
        }

        let offset_from_end = len - index - 1;
        let mut cur;
        if index <= offset_from_end {
            cur = self.head;
            for _ in 0..index {
                match cur.take() {
                    None => cur = self.head,
                    Some(node) => unsafe {
                        cur = (*node.as_ptr()).next;
                    },
                }
            }
        } else {
            cur = self.tail;
            for _ in 0..offset_from_end {
                match cur.take() {
                    None => cur = self.tail,
                    Some(node) => unsafe {
                        cur = (*node.as_ptr()).prev;
                    },
                }
            }
        }

        Ok(cur)
    }

    pub fn insert_by_index(&mut self, index: usize, val: T) -> Result<(), Box<dyn Error>> {
        let len = self.length;
        if index > len {
            return Err(Box::new(IndexOutOfRangeError));
        }
        if index == 0 {
            return Ok(self.push_front(val));
        } else if index == len {
            return Ok(self.push_back(val));
        }

        unsafe {
            let mut spliced_node = Box::new(Node::new(val));
            let before_node = self._get_by_index_mut(index - 1)?;
            // 这里使用 unwrap() 直接获取节点的值是因为，我们能够保证这些节点一定不为 None
            let after_node = before_node.unwrap().as_mut().next;
            spliced_node.next = after_node;
            spliced_node.prev = before_node;
            let spliced_node = NonNull::new(Box::into_raw(spliced_node));

            before_node.unwrap().as_mut().next = spliced_node;
            after_node.unwrap().as_mut().prev = spliced_node;
        }

        self.length += 1;

        Ok(())
    }

    // 直接将节点移除，并将在节点存放元素的所有权返回给方法调用者
    pub fn remove_by_index(&mut self, index: usize) -> Result<T, Box<dyn Error>> {
        let len = self.length;
        if index >= len {
            return Err(Box::new(IndexOutOfRangeError));
        }

        if index == 0 {
            return Ok(self.pop_front().unwrap());
        } else if index == len - 1 {
            return Ok(self.pop_back().unwrap());
        }

        let cur = self._get_by_index_mut(index)?.unwrap();
        self.unlink_node(cur);

        unsafe {
            let unlinked_node = Box::from_raw(cur.as_ptr());
            Ok(unlinked_node.val)
        }
    }

    #[inline]
    fn unlink_node(&mut self, mut node: NonNull<Node<T>>) {
        let node = unsafe { node.as_mut() };

        match node.prev {
            None => self.head = node.next,
            Some(prev) => unsafe {
                (*prev.as_ptr()).next = node.next;
            },
        }

        match node.next {
            None => self.tail = node.prev,
            Some(next) => unsafe {
                (*next.as_ptr()).prev = node.prev;
            },
        }

        self.length -= 1;
    }

    // Rust 中的范型和 C++ 的实现方式非常类似，即：
    // 对每一种具体类型生成其对应的代码，而非类似于 Java 中的类型擦除后进行类型转换，从而实现了：零成本抽象；
    // 同时，Rust 在编译时会分析究竟有哪些类型满足了范型约束，而只为那些满足了约束的具体类型实现方法！
    pub fn contains(&self, elem: &T) -> bool
    where
        T: PartialEq,
    {
        self.iter().any(|e| e == elem)
    }
}

impl<T> Default for LinkedList<T> {
    fn default() -> Self {
        Self::new()
    }
}
// 为实现Debug元素的链表实现遍历输出：traverse()
impl<T: Debug> LinkedList<T> {
    pub fn tranverse(&self) {
        print!("{{ ");
        for (index, value) in self.iter().enumerate() {
            print!(" [{}: {:?}] ", index, *value);
        }
        println!(" }}");
    }
}

// IntoIter 会获取整个链表所有节点的所有权，因此直接将链表的所有权转移至 IntoIter 中即可
impl<T> LinkedList<T> {
    pub fn into_iter(self) -> IntoIter<T> {
        IntoIter { list: self }
    }

    pub fn iter(&self) -> Iter<T> {
        Iter {
            head: self.head,
            tail: self.tail,
            len: self.length,
            _marker: PhantomData,
        }
    }

    pub fn iter_mut(&mut self) -> IterMut<T> {
        IterMut {
            head: self.head,
            tail: self.tail,
            len: self.length,
            _marker: PhantomData,
        }
    }
}

pub struct IntoIter<T> {
    list: LinkedList<T>,
}

impl<T> Iterator for IntoIter<T> {
    type Item = T;

    #[inline]
    fn next(&mut self) -> Option<Self::Item> {
        self.list.pop_front()
    }

    #[inline]
    fn size_hint(&self) -> (usize, Option<usize>) {
        (self.list.length, Some(self.list.length))
    }
}

impl<T> DoubleEndedIterator for IntoIter<T> {
    #[inline]
    fn next_back(&mut self) -> Option<Self::Item> {
        self.list.pop_back()
    }
}

// !由于 IntoIter 获取了整个链表的所有权，而我们是通过裸指针实现的链表，即我们需要手动管理这部分内存；
// !因此，我们需要手动为 IntoIter 实现 Drop Trait，以确保在 IntoIter 退出作用域后，能够准备的释放掉那些还没有被 move 出去的元素！
impl<T> Drop for IntoIter<T> {
    fn drop(&mut self) {
        // 直接通过 for 循环将 IntoIter 中还未被消费的元素直接取出来，然后忽略掉即可
        // 这里的 for _ in &mut *self {} 实际上就是调用的迭代器本身的 next 方法去取元素；
        // 而 next 是调用的链表的 pop_front 方法，该方法最终会调用 Box::from_raw 将裸指针转为具体的元素返回，因此实现了内存释放；
        for _ in &mut *self {}
        println!("IntoIter dropped");
    }
}

// 对于 Iter 和 IterMut 而言，我们需要 Copy 当前链表的头节点和尾节点，而非获取链表的所有权
// !同时，对于 Iterator 的 Item 如果是引用类型，则需要指定对应元素的生命周期
// 但是由于 head 和 tail 中存放的是裸指针（即表示，其内存分配是由我们来管理的！），
// !因此此时再次需要使用 PhantomData 来避免编译器对于生命周期的检查问题；

// T 必须活得比'a长
pub struct Iter<'a, T: 'a> {
    head: Option<NonNull<Node<T>>>,
    tail: Option<NonNull<Node<T>>>,
    len: usize,
    _marker: PhantomData<&'a Node<T>>,
}

// 相比于 IntoIter，在实现 Iter 时，我们需要自己手动维护 head 和 tail 裸指针；
// 我们不需要为特别为 Iter 实现 Drop 方法，因为 Iter 中的所有类型均已经由 Rust 标准库实现了 Drop
impl<'a, T> Iterator for Iter<'a, T> {
    type Item = &'a T;

    #[inline]
    fn next(&mut self) -> Option<Self::Item> {
        if self.len == 0 {
            None
        } else {
            self.head.map(|node| {
                self.len -= 1;
                unsafe {
                    let node = &*node.as_ptr();
                    self.head = node.next;
                    &node.val
                }
            })
        }
    }
}

impl<'a, T> DoubleEndedIterator for Iter<'a, T> {
    fn next_back(&mut self) -> Option<Self::Item> {
        if self.len == 0 {
            None
        } else {
            self.tail.map(|node| {
                self.len -= 1;
                unsafe {
                    let node = &*node.as_ptr();
                    self.tail = node.prev;
                    &node.val
                }
            })
        }
    }
}

pub struct IterMut<'a, T: 'a> {
    head: Option<NonNull<Node<T>>>,
    tail: Option<NonNull<Node<T>>>,
    len: usize,
    _marker: PhantomData<&'a mut Node<T>>,
}

// IterMut 的实现和 Iter 的实现几乎完全一致，只是将类型换为了：type Item = &'a mut T；
impl<'a, T> Iterator for IterMut<'a, T> {
    type Item = &'a mut T;

    #[inline]
    fn next(&mut self) -> Option<Self::Item> {
        if self.len == 0 {
            None
        } else {
            self.head.map(|node| {
                self.len -= 1;
                unsafe {
                    let node = &mut *node.as_ptr();
                    self.head = node.next;
                    &mut node.val
                }
            })
        }
    }
}

// 现在，我们还需要为链表本身实现 Drop Trait：
// 以确保在链表退出其作用域后（此后再也无法访问此链表），内部元素的内存能够正常的被释放；
impl<T> Drop for LinkedList<T> {
    fn drop(&mut self) {
        // 在这里，我们定义了一个 DropGuard，其内部只有 LinkedList 类型的属性，并再次为其也实现了 Drop Trait：

        struct DropGuard<'a, T>(&'a mut LinkedList<T>);
        impl<'a, T> Drop for DropGuard<'a, T> {
            fn drop(&mut self) {
                // Continue the same loop we do below. This only runs when a destructor has
                // panicked. If another one panics this will abort.
                while self.0.pop_front().is_some() {}
            }
        }

        // 此处如此设计的原因是：
        // 确保在执行下面这段释放链表元素占用内存的代码时：
        // !如果出现了 panic，则此时 DropGuard 可以再次尝试释放内存；
        while let Some(node) = self.pop_front() {
            let guard = DropGuard(self);
            drop(node);
            mem::forget(guard);
        }

        println!("LinkedList dropped");
    }
}

impl<T> LinkedList<T> {
    // 这是得益于我们为双向链表实现了 Drop Trait；
    // 因此，我们可以直接创建一个新的空双向链表来直接覆盖原链表，来实现 clear() 方法；
    // 而原链表在退出作用域之后会自动调用其 drop 方法，清空内部的节点以及对应元素，释放内存！
    pub fn clear(&mut self) {
        *self = Self::new();
    }
}

#[derive(Debug, Clone)]
pub struct IndexOutOfRangeError;

impl std::fmt::Display for IndexOutOfRangeError {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "index out of range")
    }
}

impl Error for IndexOutOfRangeError {}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        let list: LinkedList<i32> = LinkedList::new();
        let list: LinkedList<i32> = LinkedList::default();
    }
}
