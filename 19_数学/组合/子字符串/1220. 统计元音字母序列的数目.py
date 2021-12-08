# 字符串中的每个字符都应当是小写元音字母（'a', 'e', 'i', 'o', 'u'）
# 每个元音 'a' 后面都只能跟着 'e'
# 每个元音 'e' 后面只能跟着 'a' 或者是 'i'
# 每个元音 'i' 后面 不能 再跟着另一个 'i'
# 每个元音 'o' 后面只能跟着 'i' 或者是 'u'
# 每个元音 'u' 后面只能跟着 'a'


class Solution:
    def countVowelPermutation(self, n: int) -> int:
        # endswith
        a, e, i, o, u, M = 1, 1, 1, 1, 1, 10 ** 9 + 7
        for _ in range(n - 1):
            a, e, i, o, u = (e + i + u) % M, (a + i) % M, (e + o) % M, i % M, (i + o) % M

        return (a + e + i + o + u) % M


print(Solution().countVowelPermutation(n=2))
# 输出：10
# 解释：所有可能的字符串分别是："ae", "ea", "ei", "ia", "ie", "io", "iu", "oi", "ou" 和 "ua"。
#
