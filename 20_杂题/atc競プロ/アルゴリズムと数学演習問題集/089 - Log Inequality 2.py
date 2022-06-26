# loga < b*logc 是否成立
# 1<=a,b,c<=1e18 且abc均为整数
# 即判断a<c^b

a, b, c = map(int, input().split())

if c == 1:
    print('No')
    exit(0)

cur = 1
for _ in range(b):  # 注意到就算b很大 也不会运算很多次
    cur *= c
    if cur > a:
        print('Yes')
        exit(0)

print('No')

