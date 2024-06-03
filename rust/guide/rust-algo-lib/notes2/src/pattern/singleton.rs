// https://stackoverflow.com/questions/27791532/how-do-i-create-a-global-mutable-singleton
// 让你能够保证一个类只有一个实例，并提供一个访问该实例的全局节点。
// 缺点是单元测试时比较困难，因为单例模式会在整个程序中共享状态。
//
// 根据定义，Singleton 是一个全局可变对象。
// 在 Rust 中，这是一个 static mut 对象。
// 因此，为了避免各种并发问题，
// 读取或写入可变静态变量的函数或块应标记为 unsafe 。
//
// 因此，单例模式可以被认为是不安全的。
// 然而，该模式在实践中仍然被广泛使用。
// Singleton 的一个很好的读取世界示例是一个 log crate，
// 它引入了 log! 和其他 debug! 日志记录宏，
// 在设置具体的记录器实例（如 env_logger）后，
// 您可以在代码中使用这些宏。
// 正如我们所看到的， env_logger 在后台使用 log：：set_boxed_logger，
// 它有一个 unsafe 块来设置全局记录器对象。
//
// 从 Rust 1.63 开始， Mutex::new const 您可以使用全局静态 Mutex 锁，而无需延迟初始化。
// 请参阅下面的使用 Mutex 的 Singleton 示例。

// 在 Rust 中实现 Singleton 的一种纯粹安全的方法是完全不使用全局变量，并通过函数参数传递所有内容

mod safe_singleton {
    fn change(global_state: &mut u32) {
        *global_state += 1
    }

    pub fn run() {
        let mut global_state = 0u32;
        change(&mut global_state);
        change(&mut global_state);
        println!("{}", global_state);
    }
}

// 从 Rust 1.63 开始，您可以使用全局静态 Mutex 锁，而无需延迟初始化。
mod singleton_using_mutex {
    use std::sync::Mutex;

    static ARRAY: Mutex<Vec<i32>> = Mutex::new(vec![]);

    fn call() {
        ARRAY.lock().unwrap().push(1);
    }

    pub fn run() {
        call();
        call();
        call();

        let array = ARRAY.lock().unwrap();
        println!("Called {} times: {:?}", array.len(), array);
        drop(array);

        *ARRAY.lock().unwrap() = vec![3, 4, 5];
        println!("Replaced: {:?}", ARRAY.lock().unwrap());
    }
}

fn main() {
    // safe_singleton::run();
    singleton_using_mutex::run();
}
