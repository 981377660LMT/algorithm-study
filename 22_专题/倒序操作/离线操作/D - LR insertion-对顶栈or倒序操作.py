# !插入:倒序考虑问题
# nums=[0]
# s[i]="L"表示i插入i-1元素(之前插入的元素)的左边
# s[i]="R"表示i插入i-1元素(之前插入的元素)的右边
# 求最后的序列


from collections import deque
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
s = input()

# !1.对顶栈
# stack1, stack2 = [], []
# for i, char in enumerate(s):
#     if char == "L":
#         stack2.append(i)
#     else:
#         stack1.append(i)
# print(*(stack1 + [n] + stack2[::-1]), sep=" ")

# !2.倒序操作
queue = deque([n])
for i in range(n - 1, -1, -1):
    if s[i] == "L":
        queue.append(i)
    else:
        queue.appendleft(i)
print(*queue, sep=" ")
