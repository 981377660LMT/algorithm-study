# atcoder pypy3 注意事项

1. sqrt 的精度问题 sqrt 因为精度问题会导致开根计算出错
   `isqrt可以避免精度问题`
2. global 变量查找比 nonlocal 更快 优化代码时可以不把代码放在 main 函数里运行
3. 如果写动态规划 最好用 dp 数组 而且是`滚动数组`
   不要用 lru_cache 非常慢 如果必须要 dfs 需要自己写 memo
4. int 超过 longlong 之后 速度会变慢
5. 不要用 product/combinations 来代替二重/三重循环 会慢很多
6. 如果是状态有限个的判断条件 尽量按情况一一列出执行而不是用 for 循环去执行 `for 循环会慢很多`
7. 注意版本问题 math 里没有 comb 函数
8. 预处理不要用记忆化 dfs 要用 dp 数组
9. 注意除法的精度问题 合理使用 Fraction 类
10. 如果需要多次调用 self.xxx 可以先定义一个局部变量 yyy = self.xxx 再使用，减小查找的开销
