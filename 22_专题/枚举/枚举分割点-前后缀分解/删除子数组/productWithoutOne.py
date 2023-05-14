from typing import Callable, List, TypeVar

T = TypeVar("T")


def productWithoutOne(nums: List[T], e: Callable[[], T], op: Callable[[T, T], T]) -> List[T]:
    """除自身以外数组的乘积.nums数组维护区间上的幺半群."""
    n = len(nums)
    res = [e() for _ in range(n)]
    for i in range(n - 1):
        res[i + 1] = op(res[i], nums[i])
    x = e()
    for i in range(n - 1, -1, -1):
        res[i] = op(res[i], x)
        x = op(nums[i], x)
    return res
