# 如果所读的字符是字母，则将该字母写在磁带上。
# 如果所读的字符是数字（例如 d），则整个当前磁带总共会被重复写 d-1 次。
class Solution:
    def decodeAtIndex(self, s: str, k: int) -> str:
        length = 0
        for char in s:
            if char.isalpha():
                length += 1
            elif char.isnumeric():
                length *= int(char)

        # 倒着解锁
        for char in reversed(s):
            k %= length
            if k == 0 and char.isalpha():
                return char
            if char.isnumeric():
                length //= int(char)
            elif char.isalpha():
                length -= 1

        return ''


print(Solution().decodeAtIndex(s="leet2code3", k=10))
# 输出："o"
# 解释：
# 解码后的字符串为 "leetleetcodeleetleetcodeleetleetcode"。
# 字符串中的第 10 个字母是 "o"。
