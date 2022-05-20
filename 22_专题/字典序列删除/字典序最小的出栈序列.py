# 给你一个字符队列 s，长度不超过 1e5，由小写字母组成。
# 你还有一个空栈。

# 每次你可以执行下列操作之一，直到队列和栈均为空：
# - 弹出队首字符，将其入栈。
# - 弹出栈顶字符。

# 输出字典序最小的出栈序列。

# 输入 cab
# 输出 abc
# 解释：c 入栈，a 入栈，a 出栈，b 入栈，b 出栈，c 出栈

# 输入 acdb
# 输出 abdc

# https:/lcodeforces.com/problemset/submission/797/157636876
# 贪心。
# 遍历s，设当前下标为i，算出i和i后面的最小字符c，然后不断把不超过c的出栈，然后把s[i]入栈。
# 遍历结束后把栈清空。


from collections import Counter


class Solution:
    def minSequnce(self, s: str) -> str:
        stack, remain, res = [], Counter(s), []
        min_ = ord('a')
        for char in s:
            stack.append(char)
            remain[char] -= 1
            while min_ < ord('z') and remain[chr(min_)] == 0:
                min_ += 1
            while stack and ord(stack[-1]) <= min_:
                res.append(stack.pop())
        return ''.join(res)


print(Solution().minSequnce('bac'))
