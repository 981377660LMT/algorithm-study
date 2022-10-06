# binary 的一个 子序列 如果是 非空 的且没有 前导 0 （除非数字是 "0" 本身），那么它就是一个 好 的子序列。
# 1 <= binary.length <= 105

# We count the number of subsequence that ends with 0 and ends with 1.

MOD = int(1e9 + 7)

# 940. 不同的子序列 II.-one string.py

# need to handle leading zero


class Solution:
    def numberOfUniqueGoodSubsequences(self, binary: str) -> int:
        endswith = [0] * 2
        for char in binary:
            endswith[int(char)] = sum(endswith) + int(char)
        # binary 的一个 子序列 如果是 非空 的且没有 前导 0,除非数字是 "0" 本身
        return (sum(endswith) + int("0" in binary)) % MOD


print(Solution().numberOfUniqueGoodSubsequences("101"))
# 输入：binary = "101"
# 输出：5
# 解释：好的二进制子序列为 ["1", "0", "1", "10", "11", "101"] 。
# 不同的好子序列为 "0" ，"1" ，"10" ，"11" 和 "101" 。
