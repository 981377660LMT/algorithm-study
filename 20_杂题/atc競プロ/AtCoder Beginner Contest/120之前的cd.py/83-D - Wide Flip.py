# D - Wide Flip
# 01串区间翻转
# 找出最大的k 你可以翻转`长度>=k的区间` 使得反转后字符全为0
# !求最大的k
# n<=1e5

# 等价于最后变成一样
# 假设答案为k
# 那么后缀[k+1:]中任意一个元素是可以任意变的 : 翻转[1,k]和[1,k+1]
# 同理 前缀[:n-k]中任意一个元素是可以任意变的 : 翻转[n-k+1,n]和[n-k,n]
# 即
# 对于 ∀iϵ[1,n−k] 且s[i]≠s[i+1] 我们可以flip[i+1,n]
# 对于 ∀iϵ(n−k,n) 且s[i]≠s[i+1] 我们可以flip[1,i]
# !因此在每个相邻元素不相等处 k = min(max(x,n-x))


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    s = input()
    n = len(s)
    res = n
    for i in range(n - 1):
        if s[i] != s[i + 1]:
            cur = max(i + 1, n - i - 1)  # !前后缀中长度更大的那段
            res = min(res, cur)
    print(res)
