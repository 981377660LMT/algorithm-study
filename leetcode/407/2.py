from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 小红和小明在玩一个字符串元音游戏。

# 给你一个字符串 s，小红和小明将轮流参与游戏，小红 先 开始：

# 在小红的回合，她必须移除 s 中包含 奇数 个元音的任意 非空 子字符串。
# 在小明的回合，他必须移除 s 中包含 偶数 个元音的任意 非空 子字符串。
# 第一个无法在其回合内进行移除操作的玩家输掉游戏。假设小红和小明都采取 最优策略 。

# 如果小红赢得游戏，返回 true，否则返回 false。


# 英文元音字母包括：a, e, i, o, 和 u。

VOEWL = set("aeiou")


class Solution:
    def doesAliceWin(self, s: str) -> bool:
        arr = [c in VOEWL for c in s]
        vowelCount = sum(arr)
        return vowelCount > 0
        print(vowelCount)


#  s = "leetcoder"
print(Solution().doesAliceWin("leetcoder"))  # True
print(Solution().doesAliceWin("bbcd"))  # True
