# 2450. Number of Distinct Binary Strings After Applying Operations
# 01串区间翻转(k-flip),问最后能得到多少个不同的字符串

MOD = int(1e9 + 7)


class Solution:
    def countDistinctStrings(self, s: str, k: int) -> int:
        # !s中长度为k的子串有len(s) - k + 1个,每个子串有翻转与不翻转两个状态
        return pow(2, len(s) - k + 1, MOD)
