# 数组中所有数的逆元


from typing import List

MOD = int(1e9 + 7)


def allInv(nums: List[int]) -> List[int]:
    """线性求数组中所有数的逆元。nums 中不能包含 0。"""
    n = len(nums)
    res = [0] * (n + 1)
    res[0] = 1
    for i, v in enumerate(nums):
        res[i + 1] = res[i] * v % MOD
    inv = pow(res[-1], MOD - 2, MOD)
    res.pop()
    for i in range(n - 1, -1, -1):
        res[i] = res[i] * inv % MOD
        inv = inv * nums[i] % MOD
    return res


if __name__ == "__main__":
    nums = [1, 2, 3, 4, 5]
    res = allInv(nums)
    assert all((a * b) % MOD == 1 for a, b in zip(res, nums))
