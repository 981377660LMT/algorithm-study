from collections.abc import Callable


def feeder(get_next_item: Callable[[], str]) -> None:
    # Body
    pass


def async_query(
    on_success: Callable[[int], None], on_error: Callable[[int, Exception], None]
) -> None:
    pass


# 预期特定签名回调函数的框架可以用 Callable[[Arg1Type, Arg2Type,...], ReturnType] 实现类型提示。
# 纯 Callable 等价于 Callable[..., Any]，进而等价于 collections.abc.Callable 。

