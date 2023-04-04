# forSubset枚举某个状态的所有子集(子集的子集)

state = 0b1101
g1 = state
while g1 >= 0:
    if g1 == state or g1 == 0:  # 跳过空集和全集
        g1 -= 1
        continue
    g2 = g1 ^ state
    print(bin(g1)[2:], bin(g2)[2:])
    g1 = -1 if g1 == 0 else (g1 - 1) & state
