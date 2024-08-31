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
