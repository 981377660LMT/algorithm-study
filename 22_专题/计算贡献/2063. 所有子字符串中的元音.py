# https://leetcode-cn.com/problems/vowels-of-all-substrings/solution/zu-he-wen-ti-xun-huan-bian-li-python-jie-hoxq/
# 总结：对每个元音，计算他出现在多少个子串 (计算子串start*end个数)
VOWELS = {'a', 'e', 'i', 'o', 'u'}


class Solution:
    def countVowels(self, word: str) -> int:
        n = len(word)
        res = 0
        for i, char in enumerate(word):
            if char in VOWELS:
                res += (i + 1) * (n - i)
        return res


print(Solution().countVowels(word="aba"))
# 输出：6
# 解释：
# 所有子字符串是："a"、"ab"、"aba"、"b"、"ba" 和 "a" 。
# - "b" 中有 0 个元音
# - "a"、"ab"、"ba" 和 "a" 每个都有 1 个元音
# - "aba" 中有 2 个元音
# 因此，元音总数 = 0 + 1 + 1 + 1 + 1 + 2 = 6 。
