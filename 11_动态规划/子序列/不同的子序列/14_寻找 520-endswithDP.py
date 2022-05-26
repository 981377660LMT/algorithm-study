# 现在字符串中将某些部分换成问号‘？‘，再将'?'部分换成 '5'、'2' 和 '0'，
# 在所有可能的字符串中，位置不同的为 "520" 的子序列共有多少个？
# 答案会很大，所以请你模 998244353 返回
MOD = 998244353


# 如果有相同形况，最好从后向前算
class Solution:
    def findOccurrences(self, S):
        # write code here
        n5, n52, n520 = 0, 0, 0
        for char in S:
            if char == '5':
                n5 += 1
            elif char == '2':
                n52 += n5
            elif char == '0':
                n520 += n52
            elif char == '?':
                n520 += n52
                n52 += n5
                n5 += 1

        return n520 % MOD
