# random.getrandbits(k) 生成 k 位 的非负随机整数
# 若 k 很大（> 2**53），性能会受影响；通常分段生成更灵活。
# 与 random.randint(0, 2**k-1) 功能等价，但 getrandbits 无边界检查，更快、更直接。


import random


x = random.getrandbits(128)  # 生成 128 位随机整数
print(f"{x:b}")  # 二进制表示
print(f"{x:032x}")  # 固定 32 个十六进制字符
bits = [(x >> i) & 1 for i in range(128)]
print(bits)
