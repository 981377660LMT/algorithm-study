#  给你四个整数：n 、a 、b 、c ，请你设计一个算法来找出第 n 个丑数。
#  1 <= n, a, b, c <= 10^9
import math

# 提示logn二分查找


class Solution:
    def nthUglyNumber(self, n: int, a: int, b: int, c: int) -> int:
        def check(mid: int) -> bool:
            total = mid // a + mid // b + mid // c - mid // ab - mid // ac - mid // bc + mid // abc
            return total >= n

        # 计算最小公倍数
        ab = a * b // math.gcd(a, b)
        ac = a * c // math.gcd(a, c)
        bc = b * c // math.gcd(b, c)
        abc = a * bc // math.gcd(a, bc)

        left, right = 1, int(1e20)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


# 设a和b的最小公倍数为gcd_ab，
# a和c的最小公倍数为gcd_ac，
# b和c的最小公倍数为gcd_bc，
# 三者的最小公倍数为gcd_labc。
# 那么在1到x当中有多少个丑数？


# 有x// a个数可以被a整除
# 有x// b个数可以被b整除
# 有x// c个数可以被c整除
# 有x//gcd_ab个数可以同时被a和b整除
# 有x//gcd_ac个数可以同时被a和c整除
# 有x//gcd_bc个数可以同时被b和c整除
# 有x//gcd_abc个数可以同时被a、b、c整除

