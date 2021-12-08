from functools import lru_cache

# 坐标 (x1,y1) 和 (x2,y2) 之间的距离是 |x1 - x2| + |y1 - y2|。
# 2 <= word.length <= 300

# O(n^2)
class Solution:
    def minimumDistance(self, word: str) -> int:
        word_length = len(word)

        def distance(char_a: str, char_b: str) -> int:
            if not char_a or not char_b:
                return 0

            index_a = ord(char_a) - 65
            index_b = ord(char_b) - 65

            return abs(index_a // 6 - index_b // 6) + abs(index_a % 6 - index_b % 6)

        @lru_cache(None)
        def dfs(i: int, key_a: str, key_b: str) -> int:
            if i == word_length:
                return 0

            char = word[i]

            return min(
                distance(char, key_a) + dfs(i + 1, char, key_b),
                distance(char, key_b) + dfs(i + 1, key_a, char),
            )

        return dfs(0, None, None)


print(Solution().minimumDistance(word="CAKE"))
# 输出：3
# 解释：
# 使用两根手指输入 "CAKE" 的最佳方案之一是：
# 手指 1 在字母 'C' 上 -> 移动距离 = 0
# 手指 1 在字母 'A' 上 -> 移动距离 = 从字母 'C' 到字母 'A' 的距离 = 2
# 手指 2 在字母 'K' 上 -> 移动距离 = 0
# 手指 2 在字母 'E' 上 -> 移动距离 = 从字母 'K' 到字母 'E' 的距离  = 1
# 总距离 = 3
