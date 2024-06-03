// ctrl+f 查找就是用的组合模式.
// 你可以使用它将对象组合成树状结构，并且能像使用独立对象一样使用它们。
// !对于绝大多数需要生成树状结构的问题来说，组合都是非常受欢迎的解决方案。
// 组合最主要的功能是在整个树状结构上递归调用方法并对结果进行汇总。

#[allow(unused)]

mod search_core {
    pub trait ISearch {
        fn search(&self, search_str: &str);
        fn key(&self) -> &'static str;
    }

    pub struct SearchCore {
        key: &'static str,
        modules: Vec<Box<dyn ISearch>>,
    }

    impl SearchCore {
        pub fn new(key: &'static str) -> Self {
            Self {
                key,
                modules: vec![],
            }
        }

        pub fn register_module(&mut self, module: impl ISearch + 'static) {
            self.modules.push(Box::new(module));
        }

        pub fn unregister_module(&mut self, key: &'static str) {
            self.modules.retain(|m| m.key() != key);
        }
    }

    impl ISearch for SearchCore {
        fn search(&self, search_str: &str) {
            println!("Searching for {} in core", search_str);
            for module in self.modules.iter() {
                module.search(search_str);
            }
        }

        fn key(&self) -> &'static str {
            self.key
        }
    }
}

mod search_module {
    use super::search_core::ISearch;

    pub struct FileSearchModule {
        key: &'static str,
    }

    impl FileSearchModule {
        pub fn new(key: &'static str) -> Self {
            Self { key }
        }
    }

    impl ISearch for FileSearchModule {
        fn search(&self, search_str: &str) {
            println!("Searching for {} in {}", search_str, self.key);
        }

        fn key(&self) -> &'static str {
            self.key
        }
    }
}

fn main() {
    use search_core::{ISearch, SearchCore};
    use search_module::FileSearchModule;

    let mut core = SearchCore::new("core");

    let file_module1 = FileSearchModule::new("file1");
    let file_module2 = FileSearchModule::new("file2");
    let file_module3 = FileSearchModule::new("file3");

    core.register_module(file_module1);
    core.register_module(file_module2);
    core.register_module(file_module3);

    core.search("search me");

    core.unregister_module("file2");

    core.search("search me");
}
