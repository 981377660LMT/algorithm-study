# 3327. 判断 DFS 字符串是否是回文串
# https://leetcode.cn/problems/check-if-dfs-strings-are-palindromes/description/
# !python字符串拼接超级快(长度为1e5左右)


from typing import List, Tuple


class Solution:
    def findAnswer(self, parent: List[int], s: str) -> List[bool]:
        n = len(s)
        adjList = [[] for _ in range(n)]
        for i, p in enumerate(parent):
            if p != -1:
                adjList[p].append(i)

        res = [False] * n

        def dfs(cur: int) -> Tuple[str, str]:
            """返回前序遍历的字符串和后序遍历的字符串."""
            s1, s2 = "", ""
            for next_ in adjList[cur]:
                a, b = dfs(next_)
                s1 += a
                s2 = b + s2

            s1 += s[cur]
            s2 = s[cur] + s2
            res[cur] = s1 == s2
            return s1, s2

        dfs(0)
        return res
