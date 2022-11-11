"""
二进制字符串重新安排顺序需要的时间

在一秒之中，所有 子字符串 "01" 同时 被替换成 "10" 。这个过程持续进行到没有 "01" 存在。
请你返回完成这个过程所需要的秒数。
"""


class Solution:
    def secondsToRemoveOccurrences1(self, s: str) -> int:
        res = 0
        while "01" in s:
            s = s.replace("01", "10")
            res += 1
        return res

    def secondsToRemoveOccurrences2(self, s: str) -> int:
        """https://leetcode.cn/problems/time-needed-to-rearrange-a-binary-string/solution/by-newhar-o6a1/"""
        # 非常想以前单调栈那道题1_stack/单调栈/倒序遍历/2289. 使数组按非递减顺序排列-单调栈携带额外信息.py
        # 从左到右枚举每个1，每个1的最大移动次数等于max(前一个1的最大移动次数+1，前面0的总数)
        res, zero = 0, 0
        for char in s:
            if char == "0":
                zero += 1
            elif zero:
                res = max(res + 1, zero)
        return res

    def secondsToRemoveOccurrences3(self, s: str) -> int:
        """https://leetcode.cn/problems/time-needed-to-rearrange-a-binary-string/solution/by-endlesscheng-pq2x/"""
        # dp[i]表示前i个字符移动所需的秒数
        # dp[i] = dp[i-1] if s[i] == "0" else max(dp[i-1]+1,preZero)
        # 注意dp[i]会被前面的1堵住
        n = len(s)
        dp, zero = [0] * (n + 1), 0
        for i in range(1, n + 1):
            if s[i - 1] == "0":
                zero += 1
                dp[i] = dp[i - 1]
            elif zero:
                dp[i] = max(dp[i - 1] + 1, zero)
        return dp[-1]


print(Solution().secondsToRemoveOccurrences3(s="0110101"))
