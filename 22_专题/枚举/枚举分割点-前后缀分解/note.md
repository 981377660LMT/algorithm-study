## 枚举分割点+前后缀分解

问题特征：

1. 寻找两个/三个 子数组的问题
2. 问题的前后缀是对称的，枚举分割点，求左半(右闭)和右半(左开)的情况
3. **移除一段子数组**后的最值/性质 => 前缀+后缀中考虑，一般是预处理前后缀

处理前后缀，枚举分割点
`689. 三个无重叠子数组的最大和`
`1031. 两个非重叠子数组的最大和 copy`
`1477. 找两个和为目标值且不重叠的子数组-前后缀分解`

```Python
注意 accumulate 里的泛型 如果传了initial 那么 T就是initial的类型 S就是iterable的类型 (类似reduce)
class accumulate(Iterator[_T], Generic[_T]):
    if sys.version_info >= (3, 8):
        @overload
        def __init__(self, iterable: Iterable[_T], func: None = ..., *, initial: _T | None = ...) -> None: ...
        @overload
        def __init__(self, iterable: Iterable[_S], func: Callable[[_T, _S], _T], *, initial: _T | None = ...) -> None: ...
    else:
        def __init__(self, iterable: Iterable[_T], func: Callable[[_T, _T], _T] | None = ...) -> None: ...
```

不容易写错的写法:

1.  定义函数 makeDp 返回前缀 dp 数组
2.  求出前后缀
    ```python
    preDp = makeDp(XXX)
    sufDp = makeDp(XXX[::-1])[::-1]
    ```
3.  枚举分割点计算答案
    ```python
    for i in range(n):
       res+=preDp[i]*sufDp[i+1]
    ```
