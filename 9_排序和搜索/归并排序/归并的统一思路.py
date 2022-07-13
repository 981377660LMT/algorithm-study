# 将线性操作变为指数次操作


from typing import Any, TypeVar

R = TypeVar("R")


def merge(left: int, right: int) -> R:
    if left == right:
        ...
    if left + 1 == right:
        ...
    mid = (left + right) // 2
    leftRes = merge(left, mid)
    rightRes = merge(mid + 1, right)
    return api(leftRes, rightRes)


def api(left: R, right: R, /, *args: Any, **kwargs: Any) -> R:
    """里面是某个api的接口 作用于left和right上"""
    ...


if __name__ == "__main__":
    nums = [1, 2, 3, 4, 5, 6]
    print(merge(0, len(nums) - 1))
