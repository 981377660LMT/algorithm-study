# param stickers  stickers 长度范围是 [1, 50]
# param target  target 的长度在 [1, 15] 范围内，由小写字母组成
# 拼出目标 target 所需的最小贴纸数量是多少？如果任务不可能，则返回 -1。

from typing import List
from collections import Counter
from functools import lru_cache


class Solution:
    def minStickers(self, stickers: List[str], target: str) -> int:
        @lru_cache(None)
        def dfs(target: str) -> int:
            if not target:
                return 0
            res = float('inf')
            for sticker in stickers:
                if target[0] not in sticker:
                    continue
                replacedWord = addSticker(sticker, target)
                res = min(res, dfs(replacedWord) + 1)
            return res

        # 耗尽Counter删除word里的字符
        def addSticker(sticker: Counter, word: str) -> str:
            for char in sticker:
                word = word.replace(char, '', sticker[char])
            return word

        stickers = [Counter(s) for s in stickers]
        res = dfs(target)
        return res if res != float('inf') else -1


print(Solution().minStickers(["with", "example", "science"], "thehat"))
