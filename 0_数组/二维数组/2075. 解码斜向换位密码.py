class Solution:
    def decodeCiphertext(self, encodedText: str, rows: int) -> str:
        # !s[i] and s[i+cols+1] 是同组相邻字符
        cols = len(encodedText) // rows
        res = []

        # 每个分组第一个字符
        for c in range(cols):
            cur = c
            while cur < len(encodedText):
                res.append(encodedText[cur])
                cur += cols + 1

        return ''.join(res).rstrip()


print(Solution().decodeCiphertext(encodedText="ch   ie   pr", rows=3))
# 输出："cipher"
# 解释：此示例与问题描述中的例子相同。
