"""最大波动的子字符串
求 子串中出现次数 最多 的字符次数与出现次数 最少 的字符次数之差 的 最大值
"""

from itertools import combinations


class Solution:
    def largestVariance(self, s: str) -> int:
        """时间复杂度O(26*26*n)"""

        def cal(s1: str, s2: str) -> int:
            # s1 最多 s2 最少 注意必须包含 s2
            # 用一个变量记录包含s2时的差值
            res = 0
            curSum, curSumWithS2 = 0, -int(1e18)  # 一开始没有s2
            for char in s:
                if char == s1:
                    curSum += 1
                    curSumWithS2 += 1
                elif char == s2:
                    curSum -= 1
                    curSumWithS2 = curSum  # 注意这里更新最大值

                if curSum < 0:  # 全不要了
                    curSum = 0
                if curSumWithS2 > res:
                    res = curSumWithS2  # if 代替 max 2000+ms
                # res = max(res, maxSumWithS2)  # 用 max 7000+ms
            return res

        allChar = list(set(s))
        res = 0
        for s1, s2 in combinations(allChar, 2):
            res = max(res, cal(s1, s2), cal(s2, s1))
        return res

    def largestVariance2(self, s: str) -> int:
        """时间复杂度O(26*n)"""

        res = 0
        curSum = [[0] * 26 for _ in range(26)]
        curSumWiths2 = [[-int(1e18)] * 26 for _ in range(26)]
        for char in s:
            ord_ = ord(char) - ord("a")
            for i in range(26):
                if i == ord_:
                    continue
                curSum[ord_][i] += 1
                curSumWiths2[ord_][i] += 1
                curSum[i][ord_] -= 1
                curSumWiths2[i][ord_] = curSum[i][ord_]
                if curSum[i][ord_] < 0:
                    curSum[i][ord_] = 0
                res = max(res, curSumWiths2[ord_][i], curSumWiths2[i][ord_])
        return res


print(Solution().largestVariance("aababbb"))
print(Solution().largestVariance("abcde"))
print(Solution().largestVariance("bbaabbaabbaabb"))  # 2
print(Solution().largestVariance("lripaa"))  # 1
