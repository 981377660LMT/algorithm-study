# s[i] and s[i+cols+1] are adjacent characters in original text.
class Solution:
    def decodeCiphertext(self, encodedText: str, rows: int) -> str:
        cols = len(encodedText) // rows
        res = []

        for start in range(cols):
            while start < len(encodedText):
                res.append(encodedText[start])
                start += cols + 1

        return ''.join(res).rstrip()


print(Solution().decodeCiphertext(encodedText="ch   ie   pr", rows=3))
# 输出："cipher"
# 解释：此示例与问题描述中的例子相同。
