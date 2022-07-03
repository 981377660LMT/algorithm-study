# max函数非常慢 尽量用 if 比大小代替
from string import ascii_lowercase

# https://leetcode.cn/problems/substring-with-largest-variance/submissions/
print(3.5 * int(1e7) * 4 // int(1e6))


class Solution:
    def largestVariance(self, s: str) -> int:
        """时间复杂度O(26*26*n)
        
        max 换成 if 2964 ms
        """

        def cal(s1: str, s2: str) -> int:
            # s1 最多 s2 最少 注意必须包含 s2
            # 用一个变量记录包含s2时的差值
            res = 0
            maxSum, maxSumWithS2 = 0, -int(1e9)  # 一开始没有s2
            for char in s:
                if char == s1:
                    maxSum += 1
                    maxSumWithS2 += 1
                elif char == s2:
                    maxSum -= 1
                    maxSumWithS2 = maxSum  # 注意这里更新答案
                    if maxSum < 0:  # 全不要了
                        maxSum = 0
                if maxSumWithS2 > res:
                    res = maxSumWithS2
                # res = max(res, maxSumWithS2)
            return res

        res = 0
        for s1 in ascii_lowercase:
            for s2 in ascii_lowercase:
                if s1 == s2:
                    continue
                res = max(res, cal(s1, s2))

        return res

    def largestVariance2(self, s: str) -> int:
        """时间复杂度O(26*26*n)
        
        用 max 	7536 ms
        """

        def cal(s1: str, s2: str) -> int:
            # s1 最多 s2 最少 注意必须包含 s2
            # 用一个变量记录包含s2时的差值
            res = 0
            maxSum, maxSumWithS2 = 0, -int(1e9)  # 一开始没有s2
            for char in s:
                if char == s1:
                    maxSum += 1
                    maxSumWithS2 += 1
                elif char == s2:
                    maxSum -= 1
                    maxSumWithS2 = maxSum  # 注意这里更新
                    if maxSum < 0:  # 全不要了
                        maxSum = 0
                # if maxSumWithS2 > res:
                #     res = maxSumWithS2
                res = max(res, maxSumWithS2)
            return res

        res = 0
        for s1 in ascii_lowercase:
            for s2 in ascii_lowercase:
                if s1 == s2:
                    continue
                res = max(res, cal(s1, s2))

        return res
