from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# Alice 和 Bob 正在玩一个游戏。最初，Alice 有一个字符串 word = "a"。

# 给定一个正整数 k。

# 现在 Bob 会要求 Alice 执行以下操作 无限次 :

# 将 word 中的每个字符 更改 为英文字母表中的 下一个 字符来生成一个新字符串，并将其 追加 到原始的 word。
# 例如，对 "c" 进行操作生成 "cd"，对 "zb" 进行操作生成 "zbac"。

# 在执行足够多的操作后， word 中 至少 存在 k 个字符，此时返回 word 中第 k 个字符的值。


# 注意，在操作中字符 'z' 可以变成 'a'。
class Solution:
    def kthCharacter(self, k: int) -> str:
        def find(depth: int, curK: int, shift: int) -> str:
            if depth == 0:
                return chr(97 + shift % 26)
            op = 1
            length = 2**depth
            mid = length // 2
            if curK <= mid:
                return find(depth - 1, curK, shift)
            return find(depth - 1, curK - mid, shift + (op == 1))

        m = 1
        log = 0
        while m < k:
            m *= 2
            log += 1
        return find(log, k, 0)


print(Solution().kthCharacter(5))  # b
print(Solution().kthCharacter(10))  # c
