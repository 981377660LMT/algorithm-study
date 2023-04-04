# 遍历bits(非常快)

state = 0b1101
while state > 0:
    bit = (state & -state).bit_length() - 1
    print(bit)  # 0, 2, 3
    state ^= 1 << bit
