from typing import List
from functools import lru_cache
from collections import defaultdict

# 1 <= words.length <= 1000
# 1 <= words[i].length <= 1000
# words 中所有单词长度相同。

# 请你返回使用 words 构造 target 的方案数
# 在构造目标字符串的过程中，你可以按照上述规定使用 words 列表中 同一个字符串 的 多个字符 。
MOD = int(1e9 + 7)


class Solution:
    def numWays(self, words: List[str], target: str) -> int:
        m, n = len(words[0]), len(target)
        # 每个位置的储备
        counter = [defaultdict(int) for _ in range(m)]
        for word in words:
            for i, char in enumerate(word):
                counter[i][char] += 1

        @lru_cache(None)
        def dfs(targetIndex: int, wordIndex: int) -> int:
            """Return number of ways to form target[targetIndex:] w/ col wordIndex."""
            if targetIndex == n:
                return 1
            if wordIndex == m:
                return 0

            # 在这个wordIndex中 选或不选
            return (
                dfs(targetIndex, wordIndex + 1)
                + dfs(targetIndex + 1, wordIndex + 1) * counter[wordIndex][target[targetIndex]]
            )

        return dfs(0, 0) % MOD


print(Solution().numWays(words=["acca", "bbbb", "caca"], target="aba"))
# 输出：6
# 解释：总共有 6 种方法构造目标串。
# "aba" -> 下标为 0 ("acca")，下标为 1 ("bbbb")，下标为 3 ("caca")
# "aba" -> 下标为 0 ("acca")，下标为 2 ("bbbb")，下标为 3 ("caca")
# "aba" -> 下标为 0 ("acca")，下标为 1 ("bbbb")，下标为 3 ("acca")
# "aba" -> 下标为 0 ("acca")，下标为 2 ("bbbb")，下标为 3 ("acca")
# "aba" -> 下标为 1 ("caca")，下标为 2 ("bbbb")，下标为 3 ("acca")
# "aba" -> 下标为 1 ("caca")，下标为 2 ("bbbb")，下标为 3 ("caca")

