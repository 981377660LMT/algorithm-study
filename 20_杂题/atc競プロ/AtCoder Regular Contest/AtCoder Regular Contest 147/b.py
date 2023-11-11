# import sys

# sys.setrecursionlimit(int(1e6))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# MOD = 998244353
# INF = int(4e18)

# n = int(input())
# perm = list(map(int, input().split()))
# mp = {num: i for i, num in enumerate(perm)}

# # A : swap pi and pi+1
# # B : swap pi and pi+2
# res = []
# for cur in range(1, n + 1):
#     index = mp[cur]
#     # 尽可能多地交换两次
#     diff = index + 1 - cur
#     while diff >= 2:
#         pre2 = perm[index - 2]
#         mp[pre2] = index
#         perm[index - 2] = cur
#         pre2, perm[index] = perm[index], pre2
#         diff -= 2
#         res.append(("B", index - 1))
#     if diff:
#         pre1 = perm[index - 1]
#         mp[pre1] = index
#         perm[index - 1] = cur
#         pre1, perm[index] = perm[index], pre1
#         res.append(("A", index))

# print(res)


# !最小化一位交换的次数 使得perm变成升序排列
