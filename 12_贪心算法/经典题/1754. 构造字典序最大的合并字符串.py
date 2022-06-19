# 抽出最大的
# 1 <= word1.length, word2.length <= 3000
class Solution:
    def largestMerge(self, word1: str, word2: str) -> str:
        w1, w2 = list(word1), list(word2)
        i = j = 0
        sb = []

        while i < len(word1) and j < len(word2):
            # 直接比后面的字符串字典序大小
            # 优化：用后缀数组的rank数组比较字典序大小
            if w1[i:] > w2[j:]:
                sb.append(w1[i])
                i += 1
            else:
                sb.append(w2[j])
                j += 1

        sb.extend(w1[i:])
        sb.extend(w2[j:])

        return ''.join(sb)


print(Solution().largestMerge(word1="cabaa", word2="bcaaa"))
# 输出："cbcabaaaaa"
# 解释：构造字典序最大的合并字符串，可行的一种方法如下所示：
# - 从 word1 中取第一个字符：merge = "c"，word1 = "abaa"，word2 = "bcaaa"
# - 从 word2 中取第一个字符：merge = "cb"，word1 = "abaa"，word2 = "caaa"
# - 从 word2 中取第一个字符：merge = "cbc"，word1 = "abaa"，word2 = "aaa"
# - 从 word1 中取第一个字符：merge = "cbca"，word1 = "baa"，word2 = "aaa"
# - 从 word1 中取第一个字符：merge = "cbcab"，word1 = "aa"，word2 = "aaa"
# - 将 word1 和 word2 中剩下的 5 个 a 附加到 merge 的末尾。
