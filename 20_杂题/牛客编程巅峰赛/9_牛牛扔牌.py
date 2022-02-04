#
#
# @param x string字符串 字符串从前到后分别是从上到下排列的n张扑克牌
# @return string字符串
#
class Solution:
    def Orderofpoker(self, x):
        def isPrime(x: int) -> bool:
            """check a number if is a prime"""
            if x <= 1:
                return False
            for num in range(2, int(x ** 0.5) + 1):
                if x % num == 0:
                    return False
            return True

        n = len(x)
        res = []
        while len(res) < n:
            curLen = len(x)
            if isPrime(curLen // 2):
                res.append(x[:2])
                x = x[2:]
            else:
                res.append(x[-2:])
                x = x[:-2]

        return ''.join(res)
