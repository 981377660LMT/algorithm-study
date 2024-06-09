# 以借用类型为参数

事实上，你应该总是用借用类型（borrowed type）,
而不是自有数据类型的借用（borrowing the owned type）。
例如&str 而非 &String, &[T] 而非 &Vec<T>, 或者 &T 而非 &Box<T>.

# 用 format!连接字符串

优点
使用 format! 连接字符串通常更加简洁和`易于阅读`。

缺点
它通常`不是最有效`的连接字符串的方法。对一个可变的 String 类型对象进行一连串的 push 操作通常是最有效率的（尤其这个字符串已经预先分配了足够的空间）

# default()方法

实现或派生 Default 的优点是，您的类型现在可以在需要 Default 实现的地方使用，最突出的是标准库中的任何 \*or_default 函数

# 将集合视为智能指针

使用集合的 Deref 特性使其像智能指针一样，提供数据的借用或者所有权。

# zmalloc

一个经常提到的问题是初始化一个全 0 的 1K 长度的向量很慢。然而，最新的 Rust 版本针对这种情况提供了一个宏调用 zmalloc，和操作系统能返回全 0 内存的速度一样快。（真的很快）

# 向闭包传递变量

默认情况下，闭包从环境中借用捕获。或者你可以用 move 闭包来将环境的所有权全给闭包。然而，一般情况下你是想传递一部分变量到闭包中，如一些数据的拷贝、传引用或者执行一些其他操作。

这种情况应**在不同的作用域里进行变量重绑定**，可读性更好。

```rs
use std::rc::Rc;

fn main() {
    let num1 = Rc::new(1);
    let num2 = Rc::new(2);
    let num3 = Rc::new(3);
    {
        // `num1` is moved
        let num2 = num2.clone(); // `num2` is cloned
        let num3 = num3.as_ref(); // `num3` is borrowed
        move || {
            *num1 + *num2 + *num3;
        }
    };
}
```

且这样在闭包定义的时候就把哪些是复制的数据搞清楚，这样结束时无论闭包有没有消耗掉这些值，都会及早 drop 掉。

闭包能用与上下文相同的变量名来用那些复制或者 move 进来的变量。

# 关于初始化的文档（Easy doc initialization）

如果一个结构体初始化操作很复杂，当写文档的时候，可以在文档中写一个使用样例的函数。
不用每次都写初始化的部分，主要写一个以这个结构体为参数的函数的用法即可。

````RS
struct Connection {
    name: String,
    stream: TcpStream,
}

impl Connection {
    /// Sends a request over the connection.
    ///
    /// # Example
    /// ```
    /// # fn call_send(connection: Connection, request: Request) {
    /// let response = connection.send_request(request);
    /// assert!(response.is_ok());
    /// # }
    /// ```
    fn send_request(&self, request: Request) {
        // ...
    }
}

````

# 临时可变性(temporary mutability)

有的时候我们需要准备和处理一些数据(init)，当处理完之后就只会读取而不修改。这种情况可以变量重绑定将其改为不可变的。

用代码块:

```rs
let data = {
    let mut data = get_vec();
    data.sort();
    data
};

// Here `data` is immutable.
```

用变量重绑定:

```rs
let mut data = get_vec();
data.sort();
let data = data;
// Here `data` is immutable.
```

# 出错时返回消耗的参数

如果一个可出错的函数消耗（移动）一个参数，则在错误中返回该参数。

```RS
pub fn send(value: String) -> Result<(), SendError> {
    println!("using {value} in a meaningful way");
    // Simulate non-deterministic fallible action.
    use std::time::SystemTime;
    let period = SystemTime::now()
        .duration_since(SystemTime::UNIX_EPOCH)
        .unwrap();
    if period.subsec_nanos() % 2 == 1 {
        Ok(())
    } else {
        Err(SendError(value))
    }
}

pub struct SendError(String);

fn main() {
    let mut value = "imagine this is very long string".to_string();

    let success = 's: {
        // Try to send value two times.
        for _ in 0..2 {
            value = match send(value) {
                Ok(()) => break 's true,
                Err(SendError(value)) => value,
            }
        }
        false
    };

    println!("success: {success}");
}
```
