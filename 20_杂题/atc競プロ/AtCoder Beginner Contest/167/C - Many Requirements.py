# C - Many Requirements
# 满足约束条件的最大得分

from typing import List, Tuple

# !1<=nums[0]<=nums[1]<=...nums[n-1]<=m
# n<=10 m<=10 q<=50  (因为n很小 所以dfs暴力搜索)


def manyRequirements(n: int, m: int, requirements: List[Tuple[int, int, int, int]]) -> int:
    def bt(index: int, pre: int, path: List[int]) -> None:
        if index == n:
            nonlocal res
            cand = sum(score for i, j, diff, score in requirements if path[j] - path[i] == diff)
            res = cand if cand > res else res
            return
        for i in range(pre, m + 1):
            path.append(i)
            bt(index + 1, i, path)
            path.pop()

    res = 0
    bt(0, 1, [])
    return res


n, m, q = map(int, input().split())
queries = []
for _ in range(q):
    i, j, diff, score = map(int, input().split())
    i, j = i - 1, j - 1
    queries.append((i, j, diff, score))

print(manyRequirements(n, m, queries))
