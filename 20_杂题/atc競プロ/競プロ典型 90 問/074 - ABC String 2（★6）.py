# 寻找变化中的不变量
# `变化`规律:a->b b->c c->a
# 合法的操作：
# 1.si=b时 将si变为a s1-si-1产生`变化`
# 1.si=c时 将si变为b s1-si-1产生`变化`

# 问操作最多的次数
# n<=60 (暗示2^60约等于1e18)
# !答案为 2^0*d1+2^1*d2+...+2^(n-1)*d(n-1)
# !1.每次操作 这个值都会至少减少1
# !2. 每次都可以只减少1 操作第一个不是a的位置即可 即 aaababca.. 的第一个b

n = int(input())
s = input()

MAPPING = {'a': 0, 'b': 1, 'c': 2}

res = 0
for i, char in enumerate(s):
    res += (1 << i) * MAPPING[char]
print(res)

