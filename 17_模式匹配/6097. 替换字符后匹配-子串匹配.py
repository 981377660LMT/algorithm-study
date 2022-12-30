# 1 <= sub.length <= s.length <= 5000
from collections import defaultdict
from typing import List

# mappings[i] = [oldi, newi] 表示你可以将 sub 中任意数目的 oldi 字符替换为 newi 。
# sub 中每个字符 不能 被替换超过一次。
# 如果使用 mappings 替换 0 个或者若干个字符，可以将 sub 变成 s 的一个子字符串，请你返回 true，否则返回 false 。

# 枚举起点匹配子串


class Solution:
    def matchReplacement1(self, s: str, sub: str, mappings: List[List[str]]) -> bool:
        # https://leetcode.cn/problems/match-substring-after-replacement/solution/mei-ju-by-oldyan-8i75/s
        # 筛选可能的子串起始位置+压位,bitpacking O(S*T/w)

        mp1 = defaultdict(int)
        for i, c in enumerate(s):
            mp1[c] |= 1 << i
        mp2 = mp1.copy()  # 每个key在indexMap中对应的所有位置
        for x, y in mappings:
            mp2[x] |= mp1[y]
        start = mp2[sub[0]]  # 可能的子串起始位置
        for i in range(1, len(sub)):
            start &= mp2[sub[i]] >> i
        return start > 0

    def matchReplacement2(self, s: str, sub: str, mappings: List[List[str]]) -> bool:
        # 枚举起点匹配子串 O(S*T)
        mp = set((x, y) for x, y in mappings)
        for i in range(len(sub), len(s) + 1):
            if all(x == y or (x, y) in mp for x, y in zip(sub, s[i - len(sub) : i])):
                return True
        return False


print(
    Solution().matchReplacement1(
        s="fool3e7bar", sub="leet", mappings=[["e", "3"], ["t", "7"], ["t", "8"]]
    )
)
