# 把字符串 s 看作是 “abcdefghijklmnopqrstuvwxyz” 的无限环绕字符串，所以 s 看起来是这样的：

# "...zabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcd...." .
# 现在给定另一个字符串 p 。返回 s 中 唯一 的 p 的 非空子串 的数量 。

# 1 <= p.length <= 1e5
from collections import defaultdict
from itertools import pairwise


class Solution:
    def findSubstringInWraproundString(self, p: str) -> int:
        """
        每一个循环子串可以由(结尾字符，长度)唯一确定 
        因此维护以每个字符结尾的最长子串长度
        """
        nums = list(map(ord, p))
        endswith, dp = defaultdict(int, {nums[0]: 1}), 1
        for pre, cur in pairwise(nums):
            if cur - pre in (1, -25):
                dp += 1
            else:
                dp = 1
            endswith[cur] = max(endswith[cur], dp)
        return sum(endswith.values())


print(Solution().findSubstringInWraproundString("zaba"))
