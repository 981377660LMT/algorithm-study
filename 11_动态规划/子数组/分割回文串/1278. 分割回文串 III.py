from functools import lru_cache

# 首先，你可以将 s 中的部分字符修改为其他的小写英文字母。
# 接着，你需要把 s 分割成 k 个非空且不相交的子串，并且每个子串都是回文串。
# 请返回以这种方式分割字符串所需修改的最少字符数。
# 1 <= k <= s.length <= 100

# 区间dp+枚举分割点
class Solution:
    def palindromePartition(self, s: str, k: int) -> int:
        # 计算修改多少个字符能变成回文
        @lru_cache(None)
        def checkPartition(left: int, right: int) -> int:
            if left >= right:
                return 0
            return checkPartition(left + 1, right - 1) + int(s[left] != s[right])

        @lru_cache(None)
        def dfs(cur: int, parts: int) -> int:
            if parts == 1:
                return checkPartition(cur, len(s) - 1)

            res = len(s)
            for splitIndex in range(cur, len(s) - parts + 1):
                res = min(res, checkPartition(cur, splitIndex) + dfs(splitIndex + 1, parts - 1))
            return res

        return dfs(0, k)


print(Solution().palindromePartition(s="abc", k=2))
# 输出：1
# 解释：你可以把字符串分割成 "ab" 和 "c"，并修改 "ab" 中的 1 个字符，将它变成回文串。
