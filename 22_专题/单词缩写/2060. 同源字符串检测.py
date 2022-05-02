# 问s1和s2是否同源
# 1 <= s1.length, s2.length <= 40
# s1 和 s2 仅由数字 1-9 和小写英文字母组成
# s1 和 s2 中连续数字数不超过 3
from functools import lru_cache

# 时间复杂度O(NMD⋅10D)
# 其中D表示连续数字串的最长长度，本题中D=3
# https://leetcode-cn.com/problems/check-if-an-original-string-exists-given-two-encoded-strings/solution/dong-tai-gui-hua-ji-lu-ke-neng-de-chang-o87gp/


class Solution:
    def possiblyEquals(self, s1: str, s2: str) -> bool:
        """爆搜

        一般双字符串dp都会带两个长度的维度
        然后这道题需要看哪个字符串的长度比较长 再加一个diff参数的维度

        思路大致是这样的 不好调试
        """

        @lru_cache(None)
        def dfs(i: int, j: int, diff: int) -> bool:
            # 当前是s1更长，拓展s2寻求匹配
            if diff >= 0:
                if j == n2:
                    return i == n1 and diff == 0
                if s2[j].isalpha():
                    # 已遍历s1末尾字符为字母，需要完全匹配
                    if i and s1[i - 1] == s2[j]:
                        if dfs(i, j + 1, diff - 1):
                            return True
                    # 已遍历s1末尾字符为数字，消耗一个rest
                    else:
                        if dfs(i, j + 1, diff - 1):
                            return True
                # s2未到结尾，且s2[j]为数字，看从数字里消耗多少个
                else:
                    count = 0
                    while j + count < n2 and s2[j + count].isdigit():
                        curNum = int(s2[j : j + count + 1])
                        if dfs(i, j + count + 1, diff - curNum):
                            return True
                        count += 1

            # 当前是s2更长，拓展s1寻求匹配
            else:
                if i == n1:
                    return j == n2 and diff == 0
                if s1[i].isalpha():
                    # 已遍历s1末尾字符为字母，需要完全匹配
                    if j and s1[i] == s2[j - 1]:
                        if dfs(i + 1, j, diff - 1):
                            return True
                    # 已遍历s1末尾字符为数字，消耗一个rest
                    else:
                        if dfs(i + 1, j, diff - 1):
                            return True
                # s2未到结尾，且s2[j]为数字，看从数字里消耗多少个
                else:
                    count = 0
                    while i + count < n1 and s1[i + count].isdigit():
                        curNum = int(s1[i : i + count + 1])
                        if dfs(i + count + 1, j, diff - curNum):
                            return True
                        count += 1

            return False

        n1, n2 = len(s1), len(s2)
        res = dfs(0, 0, 0)
        dfs.cache_clear()
        return res


print(Solution().possiblyEquals("aab", "cad"))
# 输入：s1 = "l123e", s2 = "44"
# 输出：true
# 解释："leetcode" 可以作为原字符串
# - "leetcode"
#   -> 分割：       ["l", "e", "et", "cod", "e"]
#   -> 替换：       ["l", "1", "2",  "3",   "e"]
#   -> 连接：       "l123e"，得到 s1
# - "leetcode"
#   -> 分割：       ["leet", "code"]
#   -> 替换：       ["4",    "4"]
#   -> 连接：       "44"，得到 s2

