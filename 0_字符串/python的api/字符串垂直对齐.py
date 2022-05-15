print('asas'.ljust(20, '*'))  # 相当于padRight
print(len('asas'.ljust(20)))
print('asass   '.rstrip())  # 去除右边空格

# 字符串垂直对齐


class Solution:
    def solve(self, s):
        words = s.split()
        n = max(map(len, words))
        # ljust补全矩阵填充空白，然后取出每列，去除末尾空白
        return ["".join(col).rstrip() for col in zip(*(w.ljust(n) for w in words))]


print(Solution().solve(s="abc def ghij k"))
# ['adgk', 'beh', 'cfi', '  j']


# After splitting the string we get

# [
#   "abc",
#   "def",
#   "ghij"
#   "k"
# ]
# Then, traversing vertically we can group them into "adg", "beh", "cfi", "   j".
