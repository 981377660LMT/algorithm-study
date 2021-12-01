# 你可以对 a 和 b 的二进制表示进行位翻转操作，返回能够使按位或运算   a OR b == c  成立的最小翻转次数。

# 依次比较a、b、c的各位，
# 当c的当前位是1时，a和b的当前位只需要有1个为1即可，否则需要翻转一次
# 当c的当前位是0时，a和b必须都是0
class Solution:
    def minFlips(self, a: int, b: int, c: int) -> int:
        res = 0
        for i in range(31):
            if (c >> i) & 1:
                if (a >> i) & 1 == 0 and (b >> i) & 1 == 0:
                    res += 1
            else:
                if (a >> i) & 1 == 1:
                    res += 1
                if (b >> i) & 1 == 1:
                    res += 1
        return res


print(Solution().minFlips(a=2, b=6, c=5))
