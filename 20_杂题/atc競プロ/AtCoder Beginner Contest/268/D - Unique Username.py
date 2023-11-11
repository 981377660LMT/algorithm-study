# 给定n种字符串，请用'_'将字符串串联起来，
# 类似于a_b_c__d的形式，其中每个字符串之间隔的'_'自定，
# 再给定m个T字符，要求你构造的字符和所有的T不同。
# 最后，构造的用户名必须长度在[3,16]之间
# n<=8
# !product和combinations暴力枚举

from itertools import permutations, product
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m = map(int, input().split())
words = [input() for _ in range(n)]
ban = set(input() for _ in range(m))

wordLen = sum(len(word) for word in words)
remain = 16 - wordLen - (n - 1)  # 除去每个单词和之间的下划线'_'，还剩下多少个字符
for perm in permutations(words):
    for underline in product(range(remain + 1), repeat=n - 1):  # 枚举每个'_'的个数
        res = ""
        for i in range(n):
            res += perm[i]
            if i != n - 1:
                res += "_" * (underline[i] + 1)
        if 3 <= len(res) <= 16 and res not in ban:
            print(res)
            exit(0)
print(-1)
