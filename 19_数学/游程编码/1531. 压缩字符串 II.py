# 你需要从字符串 s 中删除最多 k 个字符，以使 s 的行程长度编码长度最小。
# 请你返回删除最多 k 个字符后，s 行程长度编码的最小长度 。

# 1 <= s.length <= 100
class Solution:
    def getLengthOfOptimalCompression(self, s: str, k: int) -> int:
        ...


print(Solution().getLengthOfOptimalCompression(s="aaabcccd", k=2))
# 输出：4
# 解释：在不删除任何内容的情况下，压缩后的字符串是 "a3bc3d" ，长度为 6 。最优的方案是删除 'b' 和 'd'，这样一来，压缩后的字符串为 "a3c3" ，长度是 4 。


# 太难了 放弃
