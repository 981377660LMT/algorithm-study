# 易混淆数 (confusing number) 在旋转180°以后，可以得到和原来不同的数，且新数字的每一位都是有效的。
class Solution:
    def confusingNumber(self, n: int) -> bool:
        d = dict(zip('01689', '01986'))
        res, s = '', str(n)
        for char in s:
            if char not in d:
                return False
            res += d[char]
        return res[::-1] != s


print(Solution().confusingNumber(89))
