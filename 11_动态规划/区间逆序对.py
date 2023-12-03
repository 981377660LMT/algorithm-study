# 区间逆序对


from typing import Callable, List


def rangeInv(nums: List[int]) -> Callable[[int, int], int]:
    """
    区间逆序对.
    时空复杂度 O(n^2) 预处理, O(1) 查询.
    """
    n = len(nums)
    dp = [[0] * (n + 1) for _ in range(n + 1)]
    for left in range(n, -1, -1):
        tmp1 = dp[left]
        tmp2 = dp[left + 1] if left + 1 <= n else []
        for right in range(left, n + 1):
            if right - left <= 1:
                continue
            tmp1[right] = (
                tmp2[right] + tmp1[right - 1] - tmp2[right - 1] + (nums[left] > nums[right - 1])
            )

    def cal(start: int, end: int) -> int:
        if start >= end:
            return 0
        return dp[start][end]

    return cal


if __name__ == "__main__":
    from random import randint

    def bruteForceInv(nums: List[int]) -> Callable[[int, int], int]:
        """
        暴力求逆序对.
        时空复杂度 O(n^2) 预处理, O(1) 查询.
        """

        def cal(start: int, end: int) -> int:
            ans = 0
            for i in range(start, end):
                for j in range(i + 1, end):
                    if nums[i] > nums[j]:
                        ans += 1
            return ans

        return cal

    def check(nums: List[int]) -> None:
        inv = rangeInv(nums)
        bfInv = bruteForceInv(nums)
        for i in range(len(nums)):
            for j in range(i + 1, len(nums) + 1):
                assert inv(i, j) == bfInv(i, j)

    for i in range(10):
        check([randint(0, 100) for _ in range(100)])

    print("Good!")
