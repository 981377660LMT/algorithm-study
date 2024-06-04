// 在不暴露复杂数据结构内部细节的情况下遍历其中所有的元素。

mod users {
    pub struct UserCollection {
        users: [&'static str; 3],
    }

    impl UserCollection {
        pub fn new() -> Self {
            Self {
                users: ["Alice", "Bob", "Charlie"],
            }
        }

        pub fn iter(&self) -> UserIterator {
            UserIterator {
                index: 0,
                collection: self,
            }
        }
    }

    pub struct UserIterator<'a> {
        index: usize,
        collection: &'a UserCollection,
    }

    impl Iterator for UserIterator<'_> {
        type Item = &'static str;

        fn next(&mut self) -> Option<Self::Item> {
            if self.index < self.collection.users.len() {
                let user = Some(self.collection.users[self.index]);
                self.index += 1;
                user
            } else {
                None
            }
        }
    }
}

fn main() {
    use users::UserCollection;
    let users = UserCollection::new();
    let mut iter = users.iter();

    println!("Users:{:?}", iter.next());
    println!("Users:{:?}", iter.next());
    println!("Users:{:?}", iter.next());
    println!("Users:{:?}", iter.next());

    users.iter().for_each(|user| println!("User: {}", user));
}
