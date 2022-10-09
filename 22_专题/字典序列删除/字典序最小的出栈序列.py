# 给你一个字符队列 s，长度不超过 1e5，由小写字母组成。
# 你还有一个空栈。
# 每次你可以执行下列操作之一，直到队列和栈均为空：
# - 弹出队首字符，将其入栈。
# - 弹出栈顶字符。
# !输出字典序最小的出栈序列/出栈序列的最小字典序。
# https:/lcodeforces.com/problemset/submission/797/157636876

# !栈模拟贪心，不大于后面最小的就pop
# 遍历结束后把栈清空。


from collections import Counter


class Solution:
    def minSequnce(self, s: str) -> str:
        """实时维护后缀的min"""
        stack, remain, res = [], Counter(s), []
        min_ = ord("a")
        for char in s:
            stack.append(char)
            remain[char] -= 1
            while min_ < ord("z") and remain[chr(min_)] == 0:
                min_ += 1
            while stack and ord(stack[-1]) <= min_:
                res.append(stack.pop())
        return "".join(res)

    def minSequnce2(self, s: str) -> str:
        """预处理后缀的min"""
        sufMin = ["|"]  # !后缀最大值
        for char in s[::-1]:
            sufMin.append(min(char, sufMin[-1]))
        sufMin = sufMin[::-1][1:]
        # print(sufMin)
        stack, res = [], []
        for i, char in enumerate(s):
            stack.append(char)
            while stack and stack[-1] <= sufMin[i]:
                res.append(stack.pop())
        return "".join(res)


print(Solution().minSequnce("bac"))
print(Solution().minSequnce("cab"))
print(Solution().minSequnce2("cab"))

# 输入 cab
# 输出 abc
# 解释：c 入栈，a 入栈，a 出栈，b 入栈，b 出栈，c 出栈
# 输入 acdb
# 输出 abdc
########################################################################################
# P1750 出栈序列
# https://www.luogu.com.cn/problem/P1750
# https://www.cnblogs.com/handwer/p/13816359.html


def minSequnce2(s: str, stackCapacity: int) -> str:
    """你需要将s的字符按顺序压入一个大小为 stackCapacity 的栈并弹出
    请输出所有出栈序列中字典序最小的序列
    n <= 1e4
    """
    ...
