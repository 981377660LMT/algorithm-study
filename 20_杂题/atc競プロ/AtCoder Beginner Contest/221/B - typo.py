# 至多进行一次邻位交换操作 能否使得两个字符串相等
# 不相等时,交换前后的位置检查 O(n)
s = list(input())
t = list(input())

for i in range(len(s)):
    if s[i] != t[i]:
        if i + 1 < len(s):
            s[i], s[i + 1] = s[i + 1], s[i]
            print("Yes" if s == t else "No")
            exit(0)
print("Yes")
