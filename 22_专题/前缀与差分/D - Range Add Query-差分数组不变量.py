# 不变量，下标mod K 的值都加起来，如果一样就是yes
# 给定一个数组和一些查询，每个查询包含一个区间[left, right]，要求判断:
# !不变量：模k的分组前缀和的差不发生变化
from typing import List, Tuple


def rangeAddQuery(nums: List[int], queries: List[Tuple[int, int]]) -> List[bool]:
    def get(mod: int, left: int, right: int) -> int:
        return preSum[mod][right + 1] - preSum[mod][left]

    preSum = [[0] * (n + 1) for _ in range(k)]  # 模k余j的分组前缀和
    for j in range(k):
        for i in range(n):
            preSum[j][i + 1] = preSum[j][i] + (nums[i] if i % k == j else 0)

    res = []
    for left, right in queries:
        sums = [get(i, left, right) for i in range(k)]
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

    res = rangeAddQuery(nums, queries)
    for v in res:
        print("Yes" if v else "No")
