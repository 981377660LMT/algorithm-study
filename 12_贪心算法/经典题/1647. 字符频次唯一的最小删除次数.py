# 如果字符串 s 中 不存在 两个不同字符 频次 相同的情况，就称 s 是 优质字符串 。
# 给你一个字符串 s，返回使 s 成为 优质字符串 需要删除的 最小 字符数。
from collections import Counter


class Solution:
    def minDeletions(self, s: str) -> int:
        freq = sorted(Counter(s).values(), reverse=True)

        res = 0
        visited = set()
        for count in freq:
            while count in visited:
                count -= 1
                res += 1
            if count > 0:
                visited.add(count)
        return res


print(Solution().minDeletions(s="aaabbbcc"))
# 输出：2
# 解释：可以删除两个 'b' , 得到优质字符串 "aaabcc" 。
# 另一种方式是删除一个 'b' 和一个 'c' ，得到优质字符串 "aaabbc" 。

