# 二进制的理解
# 操作1:加1
# 操作2:乘2
# 在120次操作内 从0变为n(n<=1e18)
# 输出操作序列

# !倒序操作


n = int(input())
res = []
while n:
    if n & 1:
        res.append("A")
        n -= 1
    else:
        res.append("B")
        n //= 2

print("".join(res[::-1]))
