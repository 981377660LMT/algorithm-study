https://atcoder.jp/contests/abc400/editorial/12642

关于浮点数 `sqrt` 计算误差的详细分析：

---

## 1. 问题背景

在竞赛或工程中，经常需要对大整数 `n` 求平方根（如判断完全平方数、二分查找等）。很多人会直接用 `sqrt(n)`，但这在大数时会有精度误差，尤其是 C++ 的 `double`、Python 的 `float`、Rust 的 `f64` 等浮点类型。

### IEEE 754 binary64（double）特性

- 有效位数（尾数/仮数部）为 53 位
- 超过 2^53 的整数无法精确表示

---

## 2. 误差来源

- 当 `n` 很大（如 1e15 ~ 1e18），`sqrt(n)` 结果会被四舍五入到最接近的可表示浮点数
- 例如，`sqrt(n)` 得到的浮点数再转回整数，可能比真实的整数平方根大 1 或小 1
- 误差可能出现在：
  - `n` 转为浮点数时
  - `sqrt(n)` 结果转为整数时

### 例子

- `n1 = (2^26 + 1)^2 - 1 = 4503599761588224`
- `n2 = (ceil(2^26.5))^2 - 1 = 9007199326062755`
- 这两个数都在 2^52 ~ 2^53 之间，double 已经不能精确表示所有整数

---

## 3. 误差的实际影响

- 直接用 `int(sqrt(n))`，有时会得到比实际小 1 或大 1 的结果
- 如果用 `m = int(sqrt(n))` 判断 `m*m <= n < (m+1)*(m+1)`，可能出错

---

## 4. 解决方法

### 方法一：整数二分

最保险的做法是用整数二分查找：

```cpp
// C++
ll l = 0, r = n;
while (l < r) {
    ll m = (l + r + 1) / 2;
    if (m * m <= n) l = m;
    else r = m - 1;
}
return l; // l*l <= n < (l+1)*(l+1)
```

### 方法二：浮点数校验

如果用 `sqrt`，可以这样修正：

```cpp
uint64_t uint64_floor_sqrt(uint64_t n) {
    if (n == 0) return 0;
    uint64_t tmp_m1 = std::sqrt(n) - 1;
    return tmp_m1 * (tmp_m1 + 2) < n ? tmp_m1 + 1 : tmp_m1;
}
```

- 先用 `sqrt` 得到近似值 `m`
- 检查 `m` 和 `m-1`，取满足 `m*m <= n < (m+1)*(m+1)` 的那个

### 方法三：Python 的 math.isqrt

Python 3.8+ 提供了 `math.isqrt(n)`，直接返回精确的整数平方根，无精度问题。

---

## 5. 总结建议

- **大整数求根，优先用整数二分或 `math.isqrt`**
- 如果用浮点 `sqrt`，一定要对结果和前后 1 个整数做校验
- 不要直接用 `int(sqrt(n))` 作为最终答案

---

## 参考代码

```python
import math

def safe_isqrt(n):
    m = int(math.sqrt(n))
    if (m + 1) * (m + 1) <= n:
        m += 1
    while m * m > n:
        m -= 1
    return m
```

---

**结论：浮点 sqrt 计算大整数时有误差，竞赛/工程中要用整数二分或 isqrt，或对结果做校验。**
