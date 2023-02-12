# 前缀和后缀的子序列匹配
# 6357. 最少得分子序列-删除子段使得子序列匹配
# 给你两个字符串 s 和 t 。

# 你可以从字符串 t 中删除任意数目的字符。
# 如果没有从字符串 t 中删除字符，那么得分为 0 ，否则：
# 令 left 为删除字符中的最小下标。
# 令 right 为删除字符中的最大下标。
# !字符串的得分为 right - left + 1 。
# !请你返回使 t 成为 s 子序列的最小得分。

# 转化为删除一段子数组的问题 => 枚举前缀left,双指针/二分看后缀需要到达的位置right
from typing import Any, List, Sequence

INF = int(1e18)


class Solution:
    def minimumScore(self, s: str, t: str) -> int:
        def makeDp(curS: Sequence[Any], curT: Sequence[Any]) -> List[int]:
            dp = [0] * (m + 1)  # t中的前i个字符匹配需要消耗s中前多少个字符
            i, j = 0, 0
            while i < len(curS) and j < len(curT):
                if curS[i] == curT[j]:
                    j += 1
                    dp[j] = i + 1
                i += 1
            for k in range(j + 1, m + 1):
                dp[k] = INF  # 匹配不了的消耗无数个字符
            return dp

        n, m = len(s), len(t)
        pre, suf = makeDp(s, t), makeDp(s[::-1], t[::-1])[::-1]

        res = m

        # !对pre每一个下标 找suf中的第一个下标使得pre[i]+suf[j]<=n 找不到就是INF
        first = [INF] * (m + 1)

        # !1.二分
        # for i, v in enumerate(pre):
        #     left, right = 0, m
        #     ok = False
        #     while left <= right:
        #         mid = (left + right) // 2
        #         if suf[mid] + v <= n:
        #             right = mid - 1
        #             ok = True
        #         else:
        #             left = mid + 1
        #     if ok:
        #         first[i] = left

        # !2.双指针
        right = 0
        for left, v in enumerate(pre):
            while right < len(suf) and suf[right] + v > n:
                right += 1
            if right < len(suf) and v + suf[right] <= n:
                first[left] = right

        for left in range(m + 1):  # 删除的左端点
            res = min(res, first[left] - left)
        return max(0, res)


print(Solution().minimumScore(s="abacaba", t="bzaa"))
# print(Solution().minimumScore(s="cde", t="xyz"))
# print(Solution().minimumScore(s="cde", t="xyz"))
# "acdedcdbabecdbebda"
# "bbecddb"
# print(Solution().minimumScore(s="acdedcdbabecdbebda", t="bbecddb"))
# "dabbbeddeabbaccecaee"
# "bcbbaabdbebecbebded"
print(Solution().minimumScore(s="dabbbeddeabbaccecaee", t="bcbbaabdbebecbebded"))
