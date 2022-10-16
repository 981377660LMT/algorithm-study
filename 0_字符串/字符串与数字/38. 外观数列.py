# 外观数列

from itertools import groupby


res = ["1"]
for _ in range(50):
    pre = res[-1]
    cur = "".join(str(len(list(group))) + char for char, group in groupby(pre))
    res.append(cur)


class Solution:
    def countAndSay(self, n: int) -> str:
        return res[n - 1]


# 1.     1
# 2.     11
# 3.     21
# 4.     1211
# 5.     111221
# 第一项是数字 1
# 描述前一项，这个数是 1 即 “ 一 个 1 ”，记作 "11"
# 描述前一项，这个数是 11 即 “ 二 个 1 ” ，记作 "21"
# 描述前一项，这个数是 21 即 “ 一 个 2 + 一 个 1 ” ，记作 "1211"
# 描述前一项，这个数是 1211 即 “ 一 个 1 + 一 个 2 + 二 个 1 ” ，记作 "111221"


# 还原外观数列
# 给定的字符串是否属于外观数列

# 1.若dp[i-1]合法 且 i+a[i]<=n，那么dp[i+a[i]]也合法
# 2.若i-a[i]>=1 且 dp[i-a[i]-1]合法，那么dp[i]也合法
# https://blog.csdn.net/m0_50269977/article/details/127282968
def isCountAndSay(s: str) -> bool:
    ...
