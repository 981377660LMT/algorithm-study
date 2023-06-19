"""
O(n+m)的子序列匹配

如果可以将小写字母插入模式串 pattern 得到待查询项 query,
那么待查询项与给定模式串匹配。可以在任何位置插入每个字符，也可以不插入字符。

!对于小写字母，“只要是能匹配的就先匹配掉”，没有任何理由可以留着这个小写字母到后面去匹配
用两个指针分别指向目标串开头和匹配串开头。
如果目标串遇到和匹配串相同的字符，两个指针同时加一。
如果目标串遇到大写字母和匹配串不同，判定失败；
否则，跳过目标串的当前小写字母。
"""


from typing import List


class Solution:
    def camelMatch(self, queries: List[str], pattern: str) -> List[bool]:
        def test(longer: str, shorter: str) -> bool:
            """longer串是否与模式串shorter驼峰匹配."""
            if not shorter:
                return True
            j = 0
            for v in longer:
                if j < len(pattern) and v == shorter[j]:
                    j += 1
                elif v.isupper():
                    return False
            return j == len(shorter)

        return [test(query, pattern) for query in queries]


assert Solution().camelMatch(
    ["FooBar", "FooBarTest", "FootBall", "FrameBuffer", "ForceFeedBack"], "FB"
) == [True, False, True, True, False]
