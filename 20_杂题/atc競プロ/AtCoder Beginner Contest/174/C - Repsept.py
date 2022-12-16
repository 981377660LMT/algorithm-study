# !7,77,777...数列中第一个k的倍数是第几项？
# 如果不存在，输出-1
# 1<=k<=1e6

# 1.倍数 <=> 模k的余数为0
# 2.鸽巢原理 模必定会出现循环
def firstK(k: int) -> int:
    mod = 0
    for i in range(k + 1):
        mod = (mod * 10 + 7) % k
        if mod == 0:
            return i + 1
    return -1


print(firstK(int(input())))
