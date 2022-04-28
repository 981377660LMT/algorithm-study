# # Python Version
# def build():
#     for i in range(0, 26):
#         if tr[0][i] == 1:
#             q.append(tr[0][i])
#     while len(q) > 0:
#         u = q[0]
#         q.pop()
#         for i in range(0, 26):
#             if tr[u][i] == 1:
#                 fail[tr[u][i]] = tr[fail[u]][i]
#                 q.append(tr[u][i])
#             else:
#                 tr[u][i] = tr[fail[u]][i]


# # Python Version
# def query(t):
#     u, res = 0, 0
#     i = 1
#     while t[i] == False:
#         u = tr[u][t[i] - ord('a')]
#         j = u
#         while j == True and e[j] != -1:
#             res += e[j]
#             e[j] = -1
#             j = fail[j]
#         i += 1
#     return res
