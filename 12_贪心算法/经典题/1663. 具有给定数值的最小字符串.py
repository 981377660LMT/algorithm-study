# 因此 a 的数值为 1 ，b 的数值为 2 ，c 的数值为 3 ，以此类推。
# 给你两个整数 n 和 k 。返回 长度 等于 n 且 数值 等于 k 的 字典序最小 的字符串。

# 先全部填a，然后增加直至满足条件
class Solution:
    def getSmallestString(self, n: int, k: int) -> str:
        res = ['a'] * n
        diff = k - n

        cursor = n - 1
        while diff > 25:
            res[cursor] = 'z'
            diff -= 25
            cursor -= 1

        res[cursor] = chr(97 + diff)

        return ''.join(res)


print(Solution().getSmallestString(n=3, k=27))
# 输出："aay"
# 解释：字符串的数值为 1 + 1 + 25 = 27，它是数值满足要求且长度等于 3 字典序最小的字符串。
