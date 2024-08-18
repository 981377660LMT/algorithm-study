# abc367-f-区间哈希(区间集合哈希)
# 给定两个数组和q次询问
# 每次询问给定两个区间[l1,r1],[l2,r2]，问两个区间内的元素经过排序后是否相同

# 事先对所有数进行一个随机映射，求和后比较即可

from collections import defaultdict
from itertools import accumulate
from random import randint
from typing import List, Tuple


def solve(
    nums1: List[int], nums2: List[int], queries: List[Tuple[int, int, int, int]]
) -> List[bool]:
    pool = defaultdict(lambda: randint(1, (1 << 61) - 1))
    nums1 = [pool[x] for x in nums1]
    nums2 = [pool[x] for x in nums2]

    preSum1 = [0] + list(accumulate(nums1))
    preSum2 = [0] + list(accumulate(nums2))
    res = [False] * len(queries)
    for i, (l1, r1, l2, r2) in enumerate(queries):
        l1 -= 1
        l2 -= 1
        res[i] = preSum1[r1] - preSum1[l1] == preSum2[r2] - preSum2[l2]
    return res


if __name__ == "__main__":
    n, m = map(int, input().split())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    queries = [tuple(map(int, input().split())) for _ in range(m)]
    res = solve(nums1, nums2, queries)  # type: ignore
    for x in res:
        print("Yes" if x else "No")
