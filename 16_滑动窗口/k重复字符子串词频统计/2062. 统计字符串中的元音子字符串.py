# 元音子字符串 是
# !`仅` 由元音（'a'、'e'、'i'、'o' 和 'u'）组成的一个子字符串，
# !且必须包含 全部五种 元音。
# 给你一个字符串 word ，统计并返回 word 中 元音子字符串的数目 。


from collections import Counter

VOWELS = set("aeiou")


class Solution:
    def countVowelSubstrings(self, word: str) -> int:
        res, left, n = 0, 0, len(word)
        counter = Counter()  # 统计元音字母出现的次数
        preLeft = 0  # 记录有效的左边界的位置
        for right in range(n):
            counter[word[right]] += 1
            while left <= right and all(counter[v] > 0 for v in VOWELS):
                counter[word[left]] -= 1
                left += 1
            if word[right] not in VOWELS:
                preLeft = right + 1
            res += max(0, left - preLeft)
        return res


assert Solution().countVowelSubstrings("cuaieuouac") == 7
