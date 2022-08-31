枚举分割点一般与`预处理`一起使用
适用于 n 在 10^5 左右的算法

特征：问题分成**两半**，枚举分割点，求左半(右闭)和右半(左开)的情况
前后缀分解

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
