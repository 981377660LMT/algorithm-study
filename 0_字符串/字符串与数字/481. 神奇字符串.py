# Kolakoski 序列 没有时间复杂度优于 O(n) 的算法
# 神奇字符串 s 仅由 '1' 和 '2' 组成，并需要遵守下面的规则：
# 神奇字符串 s 的神奇之处在于，串联字符串中 '1' 和 '2' 的连续出现次数可以生成该字符串。
# !给你一个整数 n ，返回在神奇字符串 s 的前 n 个数字中 1 的数目。

from itertools import accumulate


res = [1, 2, 2]
i = 2  # !i表示构造次数的位置
while len(res) < int(1e5):
    # digit = 1 if res[-1] == 2 else 2
    # res += [digit] * res[i]
    res.extend([res[-1] ^ 3] * res[i])  # 1 2 3 两个异或等于另一个
    i += 1

preSum = [0] + list(accumulate(res, lambda x, y: x + (y == 1)))


class Solution:
    def magicalString(self, n: int) -> int:
        return preSum[n]
