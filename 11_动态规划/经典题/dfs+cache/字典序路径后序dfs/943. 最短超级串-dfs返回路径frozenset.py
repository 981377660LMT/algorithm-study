from typing import List
from functools import lru_cache


# 1 <= words.length <= 12
# 1 <= words[i].length <= 20
# words[i] 由小写英文字母组成
# words 中的所有字符串 互不相同


@lru_cache(None)
def calWeight(u: str, v: str) -> int:
    return next(k for k in range(len(v), -1, -1) if u.endswith(v[:k]))


INF = 'a' * 1000


class Solution:
    def shortestSuperstring(self, words: List[str]) -> str:
        """找到以 words 中每个字符串作为子字符串的最短字符串。
        如果有多个有效最短字符串满足题目条件，返回其中 `任意一个` 即可。

        大部分状态压缩dp,Python都可以用cache和frozenset来实现
        """

        @lru_cache(None)
        def dfs(groups: frozenset, cur: str) -> str:
            if len(groups) == 1:
                return cur
            nextGroup = groups - {cur}
            res = INF
            for next in nextGroup:
                nPath, weight = dfs(nextGroup, next), calWeight(cur, next)
                cand = cur + nPath[weight:]
                if len(cand) < len(res):
                    res = cand
            return res

        return min((dfs(frozenset(words), start) for start in words), key=len)


print(Solution().shortestSuperstring(words=["catg", "ctaagt", "gcta", "ttca", "atgcatc"]))
# 输入：words = ["catg","ctaagt","gcta","ttca","atgcatc"]
# 输出："gctaagttcatgcatc"
