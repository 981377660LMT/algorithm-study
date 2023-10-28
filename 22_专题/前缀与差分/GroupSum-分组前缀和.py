# 环形数组前缀和/环形区间前缀和


from random import randint
from typing import Callable, List


def groupPreSum(nums: List[int], mod: int) -> Callable[[int, int, int], int]:
    """模分组前缀和."""
    preSum = [0] * (len(nums) + mod)
    for i, v in enumerate(nums):
        preSum[i + mod] = preSum[i] + v

    def cal(r: int, k: int) -> int:
        if r % mod <= k:
            return preSum[r // mod * mod + k]
        return preSum[(r + mod - 1) // mod * mod + k]

    def query(start: int, end: int, key: int) -> int:
        """区间[start,end)中下标模mod为key的元素的和."""
        if start >= end:
            return 0
        key %= mod
        return cal(end, key) - cal(start, key)

    return query


if __name__ == "__main__":
    nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    gs = groupPreSum(nums, 3)
    for _ in range(100):
        # check with bf
        left, right, key = randint(0, 10), randint(0, 10), randint(0, 2)
        if left > right:
            left, right = right, left
        sum1 = gs(left, right, key)
        sum2 = sum(nums[i] for i in range(left, right) if i % 3 == key)

        assert sum1 == sum2, (left, right, key, sum1, sum2)

    # https://atcoder.jp/contests/abc288/tasks/abc288_d
    # 不变量，下标mod K 的值都加起来，如果一样就是yes
    # 给定一个数组和一些查询，每个查询包含一个区间[left, right]，要求判断:
    # !不变量：模k的分组前缀和的差不发生变化
    from typing import List, Tuple

    def rangeAddQuery(nums: List[int], k: int, queries: List[Tuple[int, int]]) -> List[bool]:
        query = groupPreSum(nums, k)

        res = []
        for left, right in queries:
            sums = [query(left, right + 1, key) for key in range(k)]
            res.append(len(set(sums)) == 1)
        return res

    if __name__ == "__main__":
        n, k = map(int, input().split())
        nums = list(map(int, input().split()))
        queries = []

        q = int(input())
        for _ in range(q):
            left, right = map(int, input().split())
            left, right = left - 1, right - 1
            queries.append((left, right))

        res = rangeAddQuery(nums, k, queries)
        for v in res:
            print("Yes" if v else "No")
