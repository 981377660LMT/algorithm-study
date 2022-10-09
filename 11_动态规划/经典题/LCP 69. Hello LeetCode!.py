"""https://leetcode.cn/problems/rMeRt2/

# 参观者只需要集齐 helloleetcode 的 13 张字母卡片即可获得力扣纪念章。
# 在展位上有一些由字母卡片拼成的单词,words[i][j] 表示第 i 个单词的第 j 个字母。
# 你可以从这些单词中取出一些卡片,但每次拿取卡片都需要消耗游戏代币,规则如下：

# !从一个单词中取一个字母所需要的代币数量,为该字母左边和右边字母数量之积
# !可以从一个单词中多次取字母,每个字母仅可被取一次
# 例如：从 example 中取出字母 a,需要消耗代币 2*4=8,字母取出后单词变为 exmple:
# 再从中取出字母 m,需要消耗代币 2*3=6,字母取出后单词变为 exple:
# !请返回取得 helloleetcode 这些字母需要消耗代币的 最少 数量。如果无法取得,返回 -1。
# 取出字母的顺序没有要求
# !取出的所有字母`恰好`可以拼成 helloleetcode

# 1 <= words.length <= 24
# 1 <= words[i].length <= 8

总结:
# !1.dp[index][c][d][e][h][l][o][t] 表示从 index 开始,取出的各个字母的数量,一共24*960个状态
# !2.每次转移枚举需要删除哪几个位置的字母,然后计算代价,取最小值
  注意这个代价是可以预处理的、枚举子集也是可以预处理的
  转移复杂度O(2^8*8)
"""

from functools import lru_cache
from collections import deque
from typing import List, Tuple


INF = int(1e12)


@lru_cache(None)
def calCost(need: Tuple[bool, ...]) -> int:
    """计算从两端取出所需字符的代价 双指针"""
    n = len(need)
    left, right = 0, n - 1
    cost = 0
    leftMoved, rightMoved = 0, 0
    while left <= right:
        while left <= right and not need[left]:
            left += 1
        while left <= right and not need[right]:
            right -= 1
        if left > right:
            break
        leftCost = (left - leftMoved) * (n - 1 - left - rightMoved)
        rightCost = (right - leftMoved) * (n - 1 - right - rightMoved)
        if leftCost <= rightCost:
            left += 1
            cost += leftCost
            leftMoved += 1
        else:
            right -= 1
            cost += rightCost
            rightMoved += 1
    return cost


# @lru_cache(None)
# def calCost(need: Tuple[bool, ...]) -> int:
#     """计算从两端取出所需字符的代价 双指针"""
#     n = len(need)
#     left, right = 0, n - 1
#     cost = 0
#     leftMoved, rightMoved = 0, 0
#     while left <= right:
#         while left <= right and not need[left]:
#             left += 1
#         while left <= right and not need[right]:
#             right -= 1
#         if left > right:
#             break
#         leftCost = (left - leftMoved) * (n - 1 - left - rightMoved)
#         rightCost = (right - leftMoved) * (n - 1 - right - rightMoved)
#         if leftCost <= rightCost:
#             left += 1
#             cost += leftCost
#             leftMoved += 1
#         else:
#             right -= 1
#             cost += rightCost
#             rightMoved += 1
#     return cost


CDETOLH = set("cdetolh")
ORDER = {
    "c": 0,
    "d": 1,
    "e": 2,
    "t": 3,
    "o": 4,
    "l": 5,
    "h": 6,
}


class Solution:
    def Leetcode(self, words: List[str]) -> int:
        @lru_cache(None)
        def dfs(
            index: int, C: int, D: int, E: int, T: int, O: int, L: int, H: int, remain: int
        ) -> int:
            if remain == 0:
                return 0
            if index == n:
                return INF

            res = INF
            word = words[index]
            for state in states[index]:
                count = 0
                counter = [0] * 7
                select = [False] * len(word)
                for i in range(len(word)):
                    if state & (1 << i):
                        count += 1
                        counter[ORDER[word[i]]] += 1
                        select[i] = True
                if (
                    counter[0] > C
                    or counter[1] > D
                    or counter[2] > E
                    or counter[3] > T
                    or counter[4] > O
                    or counter[5] > L
                    or counter[6] > H
                ):
                    continue

                cost = calCost(tuple(select))
                cand = cost + dfs(
                    index + 1,
                    C - counter[0],
                    D - counter[1],
                    E - counter[2],
                    T - counter[3],
                    O - counter[4],
                    L - counter[5],
                    H - counter[6],
                    remain - count,
                )
                if cand < res:
                    res = cand

            return res

        n = len(words)
        states = [[] for _ in range(n)]
        for i, word in enumerate(words):
            m = len(word)
            for state in range(1 << m):
                for j in range(m):
                    if (state >> j) & 1 and word[j] not in CDETOLH:
                        break
                else:
                    states[i].append(state)

        res = dfs(0, 1, 1, 4, 1, 2, 3, 1, 13)
        dfs.cache_clear()
        return res if res != INF else -1


print(Solution().Leetcode(["hello", "leetcode"]))
print(Solution().Leetcode(words=["hold", "engineer", "cost", "level"]))
print(
    Solution().Leetcode(
        words=[
            "lkhqjztn",
            "cpoipalb",
            "hrke",
            "fveuttt",
            "conrzlm",
            "tdrohwgm",
            "odzetred",
            "jekj",
            "lh",
            "kelzwh",
        ]
    )
)

print(Solution().Leetcode(words=["ecleob", "rho", "tw", "lpl", "ebolddec"]))
