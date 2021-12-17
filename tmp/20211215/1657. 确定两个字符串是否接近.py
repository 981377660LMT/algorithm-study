from collections import Counter

# 1 <= word1.length, word2.length <= 105

# 操作 1：交换任意两个 现有 字符的位置  => set 一样，即集合无序性。
# 操作 2：交换两种字符的值 => counter(freq) 一样，即freq数组相等
class Solution:
    def closeStrings(self, word1: str, word2: str) -> bool:
        return set(word1) == set(word2) and Counter(Counter(word1).values()) == Counter(
            Counter(word2).values()
        )


print(Solution().closeStrings(word1="abc", word2="bca"))
# 输出：true
# 解释：2 次操作从 word1 获得 word2 。
# 执行操作 1："abc" -> "acb"
# 执行操作 1："acb" -> "bca"
