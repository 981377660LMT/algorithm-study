from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    s = input()

    queue = deque()

    # 1. 遍历s 如果字符为'R' 翻转queue
    # 2. 否则将字符加入queue的尾部
    # 3. 求出queue
    rev = False
    for char in s:
        if char == "R":
            rev = not rev
        else:
            if rev:
                queue.appendleft(char)
            else:
                queue.append(char)
    if rev:
        queue.reverse()

    stack = []
    for char in queue:
        if stack and char == stack[-1]:
            stack.pop()
        else:
            stack.append(char)
    print("".join(stack))
