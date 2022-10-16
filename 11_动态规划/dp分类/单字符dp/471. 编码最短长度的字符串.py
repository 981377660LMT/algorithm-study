from functools import lru_cache

# 给定一个 非空 字符串，将其编码为具有最短长度的字符串。
# 编码规则是：k[encoded_string]，其中在方括号 encoded_string 中的内容重复 k 次。
# 如果编码的过程不能使字符串缩短，则不要对其进行编码
# 1 <= s.length <= 150

# 总结：
# 判断某个字符串 ss 是否能表示为 k[子串] 可以参考 lc459。


class Solution:
    @lru_cache(None)
    def encode(self, s: str) -> str:
        res = s
        if len(s) <= 4:
            return res

        # lc459 寻找是否存在循环节
        loc = (s + s).find(s, 1, -1)

        # 有循环节
        if loc < len(s):
            res = str(len(s) // loc) + "[" + self.encode(s[:loc]) + "]"

        # 枚举分割点
        for i in range(1, len(s)):
            left = self.encode(s[:i])
            right = self.encode(s[i:])
            res = min(res, left + right, key=len)

        return res


print(Solution().encode(s="aabcaabcd"))
# 输出："2[aabc]d"
# 解释："aabc" 出现两次，因此一种答案可以是 "2[aabc]d"。
