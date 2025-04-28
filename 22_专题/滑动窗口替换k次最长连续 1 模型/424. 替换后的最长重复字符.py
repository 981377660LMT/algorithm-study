from collections import defaultdict


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def characterReplacement(self, s: str, k: int) -> int:
        """
        滑动窗口 + 计数数组：
        维护窗口 [left..right]，并跟踪窗口内出现次数最多的字符 count_max。
        当窗口大小减去 count_max 大于 k 时，说明需要的替换次数超出限制，
        则收缩左边界；否则尝试更新最大长度。
        时间 O(n)，空间 O(1)（大小固定的 26 字母计数）。
        """
        counter = defaultdict(int)
        left = 0
        countMax = 0  # 窗口内单个字符的最高出现次数
        res = 0

        for right, c in enumerate(s):
            counter[c] += 1
            # !仅需在扩大时更新 countMax，即便后面收缩时 countMax 可能略大也无妨
            countMax = max2(countMax, counter[c])
            while (right - left + 1) - countMax > k:
                counter[s[left]] -= 1
                left += 1
            res = max2(res, right - left + 1)

        return res
