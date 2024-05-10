https://fasterthanli.me/articles/a-half-hour-to-learn-rust

1. 块是表达式
2. Marker traits
3. The toilet closure

```rs
fn main() {
    countdown(3, |_| ());
}
```

之所以这样称呼，是因为 `|_| ()` 看起来像一个厕所。
