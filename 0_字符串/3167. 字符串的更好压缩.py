# 3167. 字符串的更好压缩
# https://leetcode.cn/problems/better-compression-of-string/description/


# 给定一个字符串 compressed 表示一个字符串的压缩版本。格式是一个字符后面加上其出现频率。
# 例如 "a3b1a1c2" 是字符串 "aaabacc" 的一个压缩版本。
#
# 我们在以下条件下寻求 更好的压缩：
#
# 每个字符在压缩版本中只应出现 一次。
# 字符应按 字母顺序 排列。
# 返回 compressed 的更好压缩版本。


class Solution:
    def betterCompression(self, compressed: str) -> str:
        n = len(compressed)
        left = 0
        counter = [0] * 26

        while left < n:
            c = ord(compressed[left]) - 97
            right = left + 1
            while right < n and compressed[right].isdigit():
                right += 1
            counter[c] += int(compressed[left + 1 : right])
            left = right

        return "".join(chr(p + 97) + str(counter[p]) for p in range(26) if counter[p])
