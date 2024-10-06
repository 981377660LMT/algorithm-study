# 3307. 找出第 K 个字符 II
# https://leetcode.cn/problems/find-the-k-th-character-in-string-game-ii/description/
#
# Alice 和 Bob 正在玩一个游戏。最初，Alice 有一个字符串 word = "a"。
#
# 给定一个正整数 k 和一个整数数组 operations，其中 operations[i] 表示第 i 次操作的类型。
#
# 现在 Bob 将要求 Alice 按顺序执行 所有 操作：
#
# !如果 operations[i] == 0，将 word 的一份 副本追加 到它自身。
# !如果 operations[i] == 1，将 word 中的每个字符 更改 为英文字母表中的 下一个 字符来生成一个新字符串，并将其 追加 到原始的 word。例如，对 "c" 进行操作生成 "cd"，对 "zb" 进行操作生成 "zbac"。
# 在执行所有操作后，返回 word 中第 k 个字符的值。
# 注意，在第二种类型的操作中，字符 'z' 可以变成 'a'。


from typing import List


class Solution:
    def kthCharacter(self, k: int, operations: List[int]) -> str:
        def find(depth: int, curK: int, shift: int) -> str:
            if depth == 0:
                return chr(97 + shift % 26)
            op = operations[depth - 1]
            length = 2**depth
            mid = length // 2
            if curK <= mid:
                return find(depth - 1, curK, shift)
            return find(depth - 1, curK - mid, shift + (op == 1))

        return find(len(operations), k, 0)
