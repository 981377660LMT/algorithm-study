# 问号可以删除可以替换
# 返回以a开头的最长连续递增子串
# n<=1e5


from functools import lru_cache

# O(n) Visits each character at most twice and there are n of them


class Solution:
    def solve(self, s):
        @lru_cache(None)
        def dfs(index: int, pre: int) -> int:
            if index == len(s) or pre == ord('z'):
                return 0
            res = 0
            if pre == -1:
                # 选起点'a'
                if s[index] == '?':
                    res = max(res, 1 + dfs(index + 1, ord('a')), dfs(index + 1, pre))
                elif s[index] == 'a':
                    res = max(res, 1 + dfs(index + 1, ord('a')))
            else:
                if s[index] == '?':
                    res = max(res, 1 + dfs(index + 1, pre + 1), dfs(index + 1, pre))
                elif ord(s[index]) == pre + 1:
                    res = max(res, 1 + dfs(index + 1, ord(s[index])))

            return res

        return max(dfs(i, -1) for i in range(len(s)))


print(Solution().solve(s="bca???de"))
# We can turn s into "bcabcde" and "abcde" is the longest consecutively increasing substring that starts with "a"
