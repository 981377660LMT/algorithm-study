# 头部字符串满足至少包含一个 “MT” 子序列，且以 T 结尾。
# 尾部字符串需要满足至少包含一个 “MT” 子序列，且以 M 开头
# 给出一个加密后的字符串，请你找出中间被加密的字符串 S 可能的最长表示。
n = int(input())
S = input()

left, right = 0, n - 1
while S[left] != 'M':
    left += 1
while S[left] != 'T':
    left += 1
while S[right] != 'T':
    right -= 1
while S[right] != 'M':
    right -= 1

print(S[left + 1 : right])

