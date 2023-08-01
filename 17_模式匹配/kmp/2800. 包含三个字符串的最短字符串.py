# https://leetcode.cn/problems/shortest-string-that-contains-three-strings/
# 给你三个字符串 a ，b 和 c ， 你的任务是找到长度 最短 的字符串，
# 且这三个字符串都是它的 子字符串 。
# 如果有多个这样的字符串，请你返回 字典序最小 的一个。

# 请你返回满足题目要求的字符串。


from itertools import permutations
from kmp import getNext


def minimumString(a: str, b: str, c: str) -> str:
    def maxCommon(pre: str, post: str) -> int:
        """pre的后缀和post的前缀的最大公共长度"""
        cat = post + "#" + pre
        next_ = getNext(cat)
        return next_[-1]

    res = []
    for perm in permutations([a, b, c]):
        w1, w2, w3 = perm
        if w2 not in w1:
            common1 = maxCommon(w1, w2)
            w1 = w1 + w2[common1:]
        if w3 not in w1:
            common2 = maxCommon(w1, w3)
            w1 = w1 + w3[common2:]
        res.append(w1)

    return min(res, key=lambda x: (len(x), x))


# https://www.luogu.com.cn/problem/CF25E
if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    a = input()
    b = input()
    c = input()
    print(len(minimumString(a, b, c)))
