# 長さが K である S の部分列のうち、辞書順で最小であるものを出力してください
# 1≤K≤N≤1e5

# !字典序最小的长为k的子序列


import sys


input = sys.stdin.readline

n, k = map(int, input().split())
s = input()

remove = n - k
stack = []
for char in s:
    while stack and remove and stack[-1] > char:
        stack.pop()
        remove -= 1
    stack.append(char)


print(''.join(stack[:k]))

