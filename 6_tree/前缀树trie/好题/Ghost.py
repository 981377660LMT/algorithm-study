from collections import defaultdict
from typing import List

# n ≤ 10,000 where n is the length of dictionary.
# 两人轮流报的字符串必须是某个单词前缀，谁先报到整个单词谁输
# 问第一个玩家是否胜利，在两人都最好发挥下


class Solution:
    def solve(self, words: List[str]) -> int:
        def dfs(cur) -> bool:
            if '#' in cur:
                return False
            return not any(dfs(node) for node in cur.values())  # 队手报什么都不能赢

        Trie = lambda: defaultdict(Trie)
        root = Trie()
        for w in words:
            cur = root
            for c in w:
                cur = cur[c]
            cur = cur['#']

        return any(dfs(node) for node in root.values())


print(Solution().solve(words=["ghost", "ghostbuster", "gas"]))
# Here is a sample game when dictionary is ["ghost", "ghostbuster", "gas"]:

# Player 1: "g"
# Player 2: "h"
# Player 1: "o"
# Player 2: "s"
# Player 1: "t" [loses]
# If player 2 had chosen "a" as the second letter, player 1 would still lose since they'd be forced to write the last letter "s".
