# 带权(权重为等差数列)前缀和


from typing import Callable, List


def weightedPreSum(nums: List[int], first=1, diff=1) -> Callable[[int, int], int]:
    """带权(权重为等差数列)前缀和.
    权重首项为first,公差为diff.
    前缀和为: first*a0+(first+diff)*a1+...+(first+(k-1)*diff)*ak
    """
    preSum1 = [0] * (len(nums) + 1)
    preSum2 = [0] * (len(nums) + 1)
    for i, v in enumerate(nums):
        preSum1[i + 1] = preSum1[i] + diff * v
        preSum2[i + 1] = preSum2[i] + (first + i * diff) * v

    def query(start: int, end: int) -> int:
        """区间[start,end)的带权前缀和."""
        if start >= end:
            return 0
        return preSum2[end] - preSum2[start] - start * (preSum1[end] - preSum1[start])

    return query


if __name__ == "__main__":
    import random

    def bf(nums: List[int], start=1, diff=1):
        def query(s: int, e: int) -> int:
            if s >= e:
                return 0
            return sum((start + i * diff) * v for i, v in enumerate(nums[s:e]))

        return query

    nums = [random.randint(0, 100) for _ in range(100)]
    for start in range(100):
        for diff in range(100):
            S = weightedPreSum(nums, start, diff)
            S2 = bf(nums, start, diff)
            for _ in range(100):
                left, right = random.randint(0, 100), random.randint(0, 100)
                if left > right:
                    left, right = right, left
                assert S(left, right) == S2(left, right), (
                    left,
                    right,
                    S(left, right),
                    S2(left, right),
                )

    print("ok")
