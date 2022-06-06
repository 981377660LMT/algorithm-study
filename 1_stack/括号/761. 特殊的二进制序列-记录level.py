# s是一个有效的括号字符串
# 可以交换连续有效的括号字符串
# 在任意次数的操作之后，交换后的字符串按照字典序排列的最大的结果是什么？

# 如果是满足题目要求的字符串，必然起始字符是"1"，结束字符是"0"
# 首先将字符串进行分割，分割成满足条件的子串
# 然后针对每个子串，去掉首尾字符后递归
# 最后针对已经内部排好序的子串进行排序，并且链接返回
class Solution:
    def makeLargestSpecial(self, s: str) -> str:
        if not s:
            return ''

        # 子串分割
        cands = []
        level = 0
        preStart = 0
        for i, char in enumerate(s):
            level += 1 if char == '1' else -1
            # !这是一截有效的括号字符串 注意要保证首尾字符是1和0
            if level == 0:
                cands.append('1' + self.makeLargestSpecial(s[preStart + 1 : i]) + '0')
                preStart = i + 1

        return ''.join(sorted(cands, reverse=True))


print(Solution().makeLargestSpecial(s="11011000"))
# 输出: "11100100"
# 解释:
# 将子串 "10" （在S[1]出现） 和 "1100" （在S[3]出现）进行交换。
# 这是在进行若干次操作后按字典序排列最大的结果。

