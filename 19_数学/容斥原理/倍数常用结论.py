def multiSum(upper: int, num: int) -> int:
    """[1,upper]内num的倍数的和.
    即num为首项,(upper//num)*num为末项的等差数列和.
    """
    first, last = num, (upper // num) * num
    count = 1 + (last - first) // num
    return (first + last) * count // 2


def multiCount(upper: int, num: int) -> int:
    """[1,upper]内num的倍数的个数"""
    return upper // num


# 闭区间 [0,r] 内模mod与k同余的数的个数
def modCount(right: int, k: int, mod: int) -> int:
    """区间 [0,right] 内模mod与k同余的数的个数"""
    assert 0 <= k < mod
    return (right - k + mod) // mod


if __name__ == "__main__":
    # https://leetcode.cn/problems/sum-multiples/
    # 6391. 倍数求和
    class Solution:
        def sumOfMultiples(self, n: int) -> int:
            """请你计算在 [1,n] 范围内能被 3、5、7 整除的所有整数之和。"""
            return (
                multiSum(n, 3)
                + multiSum(n, 5)
                + multiSum(n, 7)
                - multiSum(n, 15)
                - multiSum(n, 21)
                - multiSum(n, 35)
                + multiSum(n, 105)
            )
