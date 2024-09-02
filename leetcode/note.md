1. rust 字符

```rust
b'a' // u8
'a' // char
"a" // &str
"a".to_string() // String
```

2. 用 iter 并不会浪费很多内存，实际上和命令式的写法使用内存几乎一样，函数式主要怕多重递归，没递归就不怕爆栈
3. 写成返回 Iterator 的函数 producer 会复杂一些，需要用 impl Trait 来声明返回类型。
   另外，因为我们`只是借用了 nums 这个 slice，所以需要显式声明返回 impl Trait 的生命周期与 nums 一样： &'a`

   ```rust
   impl Solution {
       pub fn subsets(nums: Vec<i32>) -> Vec<Vec<i32>> {
           producer(&nums).collect()
       }
   }

   fn producer<'a>(nums: &'a [i32]) -> impl Iterator<Item = Vec<i32>> + 'a {
       (0..(1 << nums.len())).map(move |mask| { // 因为我们要返回这个闭包，所以需要用 `move` 关键字把 `nums` 捕获住
           nums.iter().enumerate().filter_map(|(i, &v)| {
               if mask & (1 << i) != 0 {
                   Some(v)
               } else {
                   None
               }
           }).collect()
       })
   }
   ```

4. 很多人都喜欢用 len 或者 is_empty 去判断长度和判空，而 Rust 不像 C++，Rust 的 vec[index]这种形式的下标访问是有边界检查的，多用迭代器，或者利用返回的 Option::None 来代替判空是比较好的。
   否则经常会出现连续两次下标检查的情况，比如像 if stack.is_empty(){ stack[index] }这样的代码实际上是等于连续判了两次空的，尽量不要写，最好是写成 if let Some(xx)=stack.get(index){...}}else{...}这种形式，只会进行一次下标检查顺带判空，同样在循环里使用迭代器依然低更好的选择。

5. string 可以从 vec 构造，由于题目保证都是 ascii 字符，所以直接上 unsafe,避免检查 utf8 有效性和重新申请内存。
   字符串可以按字节操作，rust 的 char 占四个字节，除非字符串里有 utf8 多字节字符，否则还是 u8 比 char 更快更省内存
   ```rust
   unsafe { String::from_utf8_unchecked(stack) }
   ```
6. Rust 函数返回引用的不同策略
   不能返回指向局部变量的引用。您有两种选择，要么返回值，要么使用静态变量
   https://colobu.com/2019/08/13/strategies-for-returning-references-in-rust/
   https://anyu.dev/post/%E5%A6%82%E4%BD%95%E5%9C%A8-rust-%E4%B8%AD%E8%BF%94%E5%9B%9E%E5%87%BD%E6%95%B0%E5%86%85%E5%88%9B%E5%BB%BA%E7%9A%84%E5%8F%98%E9%87%8F%E7%9A%84%E5%BC%95%E7%94%A8/

   - 模式零: 使用静态变量
   - 模式一: 返回 Owned Value（返回所有权而不是引用）
   - 模式二: 返回 Boxed Value (从函数栈移到堆)
   - **模式三：重新组织代码，将 Owned Value 移动到上面的 Scope，使用引用作为函数参数**
   - **模式四: 使用回调取代返回值**

7. 默认数字 i32，默认数组下标 usize
8. 实现 IntoIterator 的好处之一就是你的类型将适用于 Rust 的 for 循环。

---

TODO:

- 模板如何声明 Interger 泛型？
