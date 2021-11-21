from collections import defaultdict


class Solution:
    def calculateTime(self, keyboard: str, word: str) -> int:
        chr_idx = defaultdict(int)
        for index, char in enumerate(keyboard):
            chr_idx[char] = index

        res, pre = 0, 0
        for char in word:
            res += abs(chr_idx[char] - pre)
            pre = chr_idx[char]

        return res


print(Solution().calculateTime("abcdefghijklmnopqrstuvwxyz", "cba"))


# 输出：4
# 解释：
# 机械手从 0 号键移动到 2 号键来输出 'c'，又移动到 1 号键来输出 'b'，接着移动到 0 号键来输出 'a'。
# 总用时 = 2 + 1 + 1 = 4.

