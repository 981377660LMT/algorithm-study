# 你需要从字符串 s 中删除最多 k 个`相邻的字符`，以使 s 的行程长度编码长度最小。
# 请你返回删除最多 k 个字符后，s 行程长度编码的最小长度 。
from itertools import groupby

# 压缩字符串


class Solution:
    def solve(self, string: str, k: int) -> int:
        def getRLELen(x: int) -> int:
            return x if x <= 1 else len(str(x)) + 1

        n = len(string)
        if n == k:
            return 0

        # 连续长度
        left = [1] * n
        for i in range(n - 1):
            if string[i] == string[i + 1]:
                left[i + 1] = left[i] + 1

        right = [1] * n
        for i in range(n - 2, -1, -1):
            if string[i] == string[i + 1]:
                right[i] = right[i + 1] + 1

        groups = [len(list(g)) for _, g in groupby(string)]

        prefix = [0] * n
        pre = 0
        i = 0
        for g in groups:
            for j in range(1, 1 + g):
                prefix[i] = pre + getRLELen(j)
                i += 1
            pre += getRLELen(g)

        suffix = [0] * n
        pre = 0
        i = n - 1
        for g in reversed(groups):
            for j in range(1, 1 + g):
                suffix[i] = pre + getRLELen(j)
                i -= 1
            pre += getRLELen(g)

        res = min(prefix[~k], suffix[k])

        # 删除哪段子数组
        for i in range(len(string) - k - 1):
            cand = prefix[i] + suffix[i + k + 1]
            l = left[i]
            r = right[i + k + 1]

            # 删除后首尾一样，需要加上
            if string[i] == string[i + k + 1]:
                cand -= getRLELen(l) + getRLELen(r)
                cand += getRLELen(l + r)
            res = min(res, cand)

        return res


print(Solution().solve(string="aaaaabbaaaaaccaaa", k=2))
