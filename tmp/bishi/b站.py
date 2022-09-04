# import sys

# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# MOD = 998244353
# INF = int(4e18)


# def calDist(word1: str, word2: str) -> int:
#     """编辑距离O(n^2)"""
#     n1, n2 = len(word1), len(word2)
#     dp = [[INF] * (n2 + 1) for _ in range(n1 + 1)]
#     dp[0][0] = 0
#     for i in range(n1):
#         for j in range(n2):
#             if word1[i] == word2[j]:
#                 dp[i + 1][j + 1] = dp[i][j]
#             else:
#                 dp[i + 1][j + 1] = min(dp[i][j + 1] + 1, dp[i + 1][j] + 1, dp[i][j] + 1)
#     return dp[-1][-1]


# s1, s2 = input().split()
# print(calDist(s1, s2))


s = input()
n = len(s)
ok = [True] * n
S = []
lmp = {"(": ")", "[": "]", "{": "}"}
rmp = {")": "(", "]": "[", "}": "{"}
for i, char in enumerate(s):
    if char in rmp:
        if not S or s[S[-1]] != rmp[char]:
            ok[i] = False
        else:
            S.pop()
    else:
        S.append(i)

for i in S:
    ok[i] = False


for i, flag in enumerate(ok):
    if not flag:
        if s[i] in rmp:
            print(f"{rmp[s[i]]},{i+1}")
        else:
            print(f"{lmp[s[i]]},{i+2}")
