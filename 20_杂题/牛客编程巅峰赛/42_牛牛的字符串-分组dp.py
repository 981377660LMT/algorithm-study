#
#
# @param s string字符串 s.size() <= 1e5
# @param k int整型 k <= s.size()
# @return int整型
# 可以选择一个位置 i 并在 i 和 i + K 处交换字符
# 新形成的字符串应字典序大于旧字符串。为了尽可能交换尽量多的步数。最多可以交换多少步呢。
# 以k为间隔 分组dp
#
class Solution:
    def turn(self, s: str, k: int):
        # write code here
        res = 0
        for start in range(min(k, len(s))):
            counter = [0] * 26
            for i in range(start, len(s), k):
                num = ord(s[i]) - ord('a')
                counter[num] += 1
                res += sum(counter[:num])
        return res
