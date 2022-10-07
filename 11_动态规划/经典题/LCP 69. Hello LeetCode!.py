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
from collections import Counter
from typing import List, Tuple


INF = int(1e18)


HELLOLEETCODE = "helloleetcode"
Good = set(list(HELLOLEETCODE))


@lru_cache(None)
def calCost(need: Tuple[bool, ...]) -> int:
    """计算取出所有字符的代价

    贪心推公式+双指针模拟
    """
    ...


class Solution:
    def Leetcode(self, words: List[str]) -> int:
        need = Counter(HELLOLEETCODE)
        cur = Counter()
        for word in words:
            for char in word:
                cur[char] += 1
        if cur & need != need:
            return -1

        # dp[index][H][E][L][O][T][C][D]
        @lru_cache(None)
        def dfs(
            index: int, H: int, E: int, L: int, O: int, T: int, C: int, D: int, remain: int
        ) -> int:
            if remain == 0:
                return 0

            # print(index, H, E, L, O, T, C, D, remain)

            if index == n:
                return INF

            # 从一个单词中取一个字母所需要的代币数量,为该字母左边和右边字母数量之积
            # 取出的所有字母`恰好`可以拼成 helloleetcode

            # 枚举取出的字母 2^8
            # 先取两头 再取中间
            # 从左边取
            res = INF
            word = words[index]
            len_ = len(word)
            left, right = 0, len_ - 1
            for state in range(1 << len_):
                todoLen = 0
                todo = [False] * len_
                flag = True  # 检查是否都在Good里面
                for i in range(len_):
                    if state & (1 << i):
                        if word[i] not in Good:
                            flag = False
                            break
                        todo[i] = True
                        todoLen += 1

                if not flag:
                    continue

                counter = dict()  # 检查是否不超过代取的
                for i in range(len_):
                    if todo[i]:
                        counter[word[i]] = counter.get(word[i], 0) + 1
                noMoreThan = True
                for key, value in counter.items():
                    if key == "h" and value > H:
                        noMoreThan = False
                        break
                    elif key == "e" and value > E:
                        noMoreThan = False
                        break
                    elif key == "l" and value > L:
                        noMoreThan = False
                        break
                    elif key == "o" and value > O:
                        noMoreThan = False
                        break
                    elif key == "t" and value > T:
                        noMoreThan = False
                        break
                    elif key == "c" and value > C:
                        noMoreThan = False
                        break
                    elif key == "d" and value > D:
                        noMoreThan = False
                        break
                if not noMoreThan:
                    continue

                # !计算花费 找到左右近的然后删除
                # !用deque模拟比双指针更好
                left, right = 0, len_ - 1
                alive = [1] * len_
                cost = 0
                while left <= right:
                    while left < right and not todo[left]:
                        left += 1
                    while left < right and not todo[right]:
                        right -= 1

                    leftLeft, rightRight = INF, INF
                    leftRight, rightLeft = 0, 0
                    if todo[left]:
                        leftLeft = sum(alive[:left])
                        leftRight = sum(alive[left + 1 :])
                    if todo[right]:
                        rightLeft = sum(alive[:right])
                        rightRight = sum(alive[right + 1 :])

                    # 哪个离哪端近 就先删谁
                    if leftLeft < rightRight:
                        cost += leftLeft * leftRight
                        alive[left] = 0
                        left += 1
                    else:
                        # print(1, todo, left, right)
                        cost += rightLeft * rightRight
                        alive[right] = 0
                        right -= 1

                cand = cost + dfs(
                    index + 1,
                    H - counter.get("h", 0),
                    E - counter.get("e", 0),
                    L - counter.get("l", 0),
                    O - counter.get("o", 0),
                    T - counter.get("t", 0),
                    C - counter.get("c", 0),
                    D - counter.get("d", 0),
                    remain - todoLen,
                )

                if cand < res:
                    res = cand
            return res

        n = len(words)
        res = dfs(
            0,
            need["h"],
            need["e"],
            need["l"],
            need["o"],
            need["t"],
            need["c"],
            need["d"],
            len(HELLOLEETCODE),
        )
        dfs.cache_clear()
        return res


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
