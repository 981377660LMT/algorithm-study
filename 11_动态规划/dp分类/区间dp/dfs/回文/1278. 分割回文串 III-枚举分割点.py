from functools import lru_cache

# 首先，你可以将 s 中的部分字符修改为其他的小写英文字母。
# 接着，你需要把 s 分割成 k 个非空且不相交的子串，并且每个子串都是回文串。
# 请返回以这种方式分割字符串所需修改的最少字符数。
# 1 <= k <= s.length <= 100


# 总结：
# 枚举分割点+记忆化dfs


class Solution:
    def palindromePartition(self, s: str, k: int) -> int:
        @lru_cache(None)
        def cal(left: int, right: int) -> int:
            """计算[left,right]修改多少个字符能变成回文"""
            if left >= right:
                return 0
            return cal(left + 1, right - 1) + int(s[left] != s[right])

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if index >= n:
                return int(1e20)
            if remain == 1:
                return cal(index, n - 1)

            res = n
            for mid in range(index, n):
                res = min(res, cal(index, mid) + dfs(mid + 1, remain - 1))
            return res

        n = len(s)
        return dfs(0, k)


print(Solution().palindromePartition(s="abc", k=2))
# 输出：1
# 解释：你可以把字符串分割成 "ab" 和 "c"，并修改 "ab" 中的 1 个字符，将它变成回文串。
