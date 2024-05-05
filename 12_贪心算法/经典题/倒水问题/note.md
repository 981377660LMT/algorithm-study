# 倒水问题结论：

给定一个数组，每次可以选择两个正数同时减 1，问做多能操作多少次.

- 一种是有一个数特别大，最后剩下这个正数.
- 否则，一定可以两两匹配到只剩一个 1，或都为 0.

```py
max_ = max(a, b, c)
sum_ = a + b + c
restSum = sum_ - max_
if max_ > restSum:
    return restSum
else:
    return sum_ // 2
```
