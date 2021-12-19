import collections


class Solution:
    def minSteps(self, s: str, t: str) -> int:
        sc = collections.Counter(s)
        tc = collections.Counter(t)

        # 缺少的字符数
        return len(list((sc - tc).elements()))


print(Solution().minSteps(s="bab", t="aba"))
# 输出：1
# 提示：用 'b' 替换 t 中的第一个 'a'，t = "bba" 是 s 的一个字母异位词。
