MOD = int(1e9 + 7)


# 给定一个字符串 s，计算 s 的 不同非空子序列 的个数
# s 仅由小写英文字母组成
# O(26n)

# 总结：
# 统计以当前字符结束的子序列有多少个
class Solution:
    def distinctSubseqII(self, s: str) -> int:
        # endswith[i] to count how many sub sequence that ends with ith character.
        endswith = [0] * 26
        for char in s:
            endswith[ord(char) - ord('a')] = sum(endswith) + 1
            print(endswith)
        return sum(endswith) % (int(1e9 + 7))


print(Solution().distinctSubseqII(s="aba"))
# 输出：6
# 解释：6 个不同的子序列分别是 "a", "b", "ab", "ba", "aa" 以及 "aba"。
print(Solution().distinctSubseqII(s="abc"))
