use std::collections::hash_map::DefaultHasher;
use std::hash::Hasher;
use std::vec::Vec;

struct MyHashMap {
    data: Vec<Vec<(i32, i32)>>,
    capacity: usize,
}

/**
 * `&self` means the method takes an immutable reference.
 * If you need a mutable reference, change it to `&mut self` instead.
 */
impl MyHashMap {
    fn new() -> Self {
        MyHashMap {
            data: vec![vec![]; 10000],
            capacity: 1024,
        }
    }

    fn put(&mut self, key: i32, value: i32) {
        let hash = self.hash(key);
        // 如果 key 已经存在，则更新 value
        for i in 0..self.data[hash].len() {
            if self.data[hash][i].0 == key {
                self.data[hash][i].1 = value;
                return;
            }
        }
        // 如果 key 不存在，则插入
        self.data[hash].push((key, value));
    }

    fn get(&self, key: i32) -> i32 {
        let hash = self.hash(key);
        for i in 0..self.data[hash].len() {
            if self.data[hash][i].0 == key {
                return self.data[hash][i].1;
            }
        }
        -1
    }

    fn remove(&mut self, key: i32) {
        let hash = self.hash(key);
        for i in 0..self.data[hash].len() {
            if self.data[hash][i].0 == key {
                self.data[hash].remove(i);
                return;
            }
        }
    }

    fn hash(&self, key: i32) -> usize {
        let mut hasher = DefaultHasher::new();
        hasher.write_i32(key);
        (hasher.finish() & 1023) as usize
    }
}
