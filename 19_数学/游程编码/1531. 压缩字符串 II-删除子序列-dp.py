# 你需要从字符串 s 中删除最多 k 个字符，以使 s 的行程长度编码长度最小。
# 请你返回删除最多 k 个字符后，s 行程长度编码的最小长度 。

# 1 <= s.length <= 100
from functools import lru_cache

# O(n^3)解法：两个状态+内部一层遍历

# to jump or not to jump


class Solution:
    def getLengthOfOptimalCompression(self, s: str, k: int) -> int:
        def getLRELen(count: int) -> int:
            return count if count <= 1 else len(str(count)) + 1

        @lru_cache(None)
        def dfs(index: int, count: int) -> int:
            if count > target:
                return int(1e20)
            if index == len(s):
                return 0 if count == target else int(1e20)

            res = int(1e20)
            sameCount = 0
            # 我们可以从当前的位置 index 开始向后遍历，只要发现后面有字符和 s[p] 相等，则选取。这样我们可以枚举选取的字符数量，进行状态转移。
            for next in range(index, len(s)):
                sameCount += int(s[next] == s[index])
                res = min(res, dfs(next + 1, count + sameCount) + getLRELen(sameCount))

            res = min(res, dfs(index + 1, count))
            return res

        target = len(s) - k
        res = dfs(0, 0)
        dfs.cache_clear()
        return res


print(Solution().getLengthOfOptimalCompression(s="aaabcccd", k=2))
# 输出：4
# 解释：在不删除任何内容的情况下，压缩后的字符串是 "a3bc3d" ，长度为 6 。最优的方案是删除 'b' 和 'd'，这样一来，压缩后的字符串为 "a3c3" ，长度是 4 。

