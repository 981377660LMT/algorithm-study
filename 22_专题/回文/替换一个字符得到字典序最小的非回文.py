class Solution:
    def solve(self, s):
        # 遍历回文只需要到 len(s)//2 -1
        for i in range(len(s) // 2):
            if s[i] != 'a':
                return s[:i] + 'a' + s[i + 1 :]
        # 全是a或者aaaaXaaaa
        return s[:-1] + 'b'

