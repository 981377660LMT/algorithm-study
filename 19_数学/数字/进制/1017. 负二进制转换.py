# 给出数字 N，返回由若干 "0" 和 "1"组成的字符串，该字符串为 N 的负二进制（base -2）表示。
# 除非字符串就是 "0"，否则返回的字符串中不能含有前导零。


class Solution:
    def baseNeg2(self, x: int):
        res = []
        while x:
            div, mod = divmod(x, 2,)
            x = -div
            res.append(mod)
        return "".join(map(str, res[::-1] or [0]))

    def base21(self, x: int):
        res = []
        while x:
            res.append(x & 1)
            x = x >> 1
        return "".join(map(str, res[::-1] or [0]))

    def baseNeg21(self, x: int):
        res = []
        while x:
            res.append(x & 1)
            x = -(x >> 1)
        return "".join(map(str, res[::-1] or [0]))

    def base22(self, N):
        if N == 0 or N == 1:
            return str(N)
        return self.base22(N >> 1) + str(N & 1)

    def baseNeg22(self, N):
        if N == 0 or N == 1:
            return str(N)
        return self.baseNeg2(-(N >> 1)) + str(N & 1)


print(Solution().baseNeg2(2))
# 输入：2
# 输出："110"
# 解释：(-2) ^ 2 + (-2) ^ 1 = 2
