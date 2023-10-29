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


def productWithoutOne2D(
    nums: List[List[T]], e: Callable[[], T], op: Callable[[T, T], T]
) -> List[List[T]]:
    """除自身以外数组的乘积二维版本."""
    row = len(nums)
    col = len(nums[0])

    unit = e()
    res = [[unit] * col for _ in range(row)]

    for i in range(row - 1, -1, -1):
        tmp1 = res[i]
        tmp2 = nums[i]
        for j in range(col - 1, -1, -1):
            tmp1[j] = unit
            unit = op(tmp2[j], unit)

    unit = e()
    for i in range(row):
        tmp1 = res[i]
        tmp2 = nums[i]
        for j in range(col):
            tmp1[j] = op(unit, tmp1[j])
            unit = op(unit, tmp2[j])

    return res


if __name__ == "__main__":
    # https://leetcode.cn/problems/construct-product-matrix/description/
    class Solution:
        def constructProductMatrix(self, grid: List[List[int]]) -> List[List[int]]:
            return productWithoutOne2D(grid, lambda: 1, lambda a, b: a * b % 12345)
