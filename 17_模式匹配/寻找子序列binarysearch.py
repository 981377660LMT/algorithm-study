# 是否存在子序列binarysearch 索引间隔相等


class Solution:
    def solve(self, s):
        b = [i for i, c in enumerate(s) if c == "b"]
        i = [i for i, c in enumerate(s) if c == "i"]
        for j in b:
            for k in i:
                if k > j:
                    if "binarysearch" in s[j :: k - j]:
                        return True
        return False
