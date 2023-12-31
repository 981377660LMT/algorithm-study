from itertools import accumulate
from random import randint
from typing import List, Tuple, Optional
from collections import defaultdict, Counter

# 给你一个长度为 偶数 n ，下标从 0 开始的字符串 s 。

# 同时给你一个下标从 0 开始的二维整数数组 queries ，其中 queries[i] = [ai, bi, ci, di] 。

# 对于每个查询 i ，你需要执行以下操作：

# 将下标在范围 0 <= ai <= bi < n / 2 内的 子字符串 s[ai:bi] 中的字符重新排列。
# 将下标在范围 n / 2 <= ci <= di < n 内的 子字符串 s[ci:di] 中的字符重新排列。
# 对于每个查询，你的任务是判断执行操作后能否让 s 变成一个 回文串 。

# 每个查询与其他查询都是 独立的 。

# 请你返回一个下标从 0 开始的数组 answer ，如果第 i 个查询执行操作后，可以将 s 变为一个回文串，那么 answer[i] = true，否则为 false 。

# 子字符串 指的是一个字符串中一段连续的字符序列。
# s[x:y] 表示 s 中从下标 x 到 y 且两个端点 都包含 的子字符串。


MOD = int(1e9 + 7)
INF = int(1e20)

pool = defaultdict(lambda: randint(1, (1 << 61) - 1))


# TODO: 遍历两个区间数组的func
class Solution:
    def canMakePalindromeQueries(self, s: str, queries: List[List[int]]) -> List[bool]:
        s1 = s[: len(s) // 2]
        s2 = s[len(s) // 2 :][::-1]
        if Counter(s1) != Counter(s2):
            return [False] * len(queries)

        n = len(s) // 2
        preSumHash1 = [0] * (n + 1)
        preSumHash2 = [0] * (n + 1)
        for i in range(n):
            preSumHash1[i + 1] = preSumHash1[i] + pool[s1[i]]
            preSumHash2[i + 1] = preSumHash2[i] + pool[s2[i]]

        isBad = [int(a != b) for a, b in zip(s1, s2)]
        badPreSum = [0] + list(accumulate(isBad))

        res = []
        for start1, end1, start2, end2 in queries:
            start2, end2 = 2 * n - end2 - 1, 2 * n - start2 - 1
            end1 += 1
            end2 += 1
            okHash1 = preSumHash1[end1] - preSumHash1[start1]
            okHash2 = preSumHash2[end2] - preSumHash2[start2]

            inters = sorted(set([0] + [start1, end1, start2, end2] + [n]))

            for a, b in zip(inters, inters[1:]):
                # 覆盖，覆盖
                print(a, b)
                ok1 = start1 <= a <= b <= end1
                ok2 = start2 <= a <= b <= end2
                if ok1 and ok2:
                    ...
                elif ok1:
                    hash2_ = preSumHash2[b] - preSumHash2[a]
                    okHash1 -= hash2_
                elif ok2:
                    hash1_ = preSumHash1[b] - preSumHash1[a]
                    okHash2 -= hash1_  # 这里需要
                else:
                    badCount = badPreSum[b] - badPreSum[a]
                    if badCount > 0:
                        res.append(False)
                        break
            else:
                res.append(okHash1 == okHash2)

        return res


# print(Solution().canMakePalindromeQueries(s="abcabc", queries=[[1, 1, 3, 5], [0, 2, 5, 5]]))
# # s = "abbcdecbba", queries = [[0,2,7,9]]
# print(Solution().canMakePalindromeQueries(s="abbcdecbba", queries=[[0, 2, 7, 9]]))
# s = "acbcab", queries = [[1,2,4,5]]
# print(Solution().canMakePalindromeQueries(s="acbcab", queries=[[1, 2, 4, 5]]))
# "odaxusaweuasuoeudxwa"
# [[0,5,10,14]]
print(Solution().canMakePalindromeQueries(s="odaxusaweuasuoeudxwa", queries=[[0, 5, 10, 14]]))
