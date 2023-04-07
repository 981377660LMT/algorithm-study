def genSubsetOfStateAscending(state: int):
    """升序枚举state所有子集的子集.
    0b1101 -> 0,1,4,5,8,9,12,13.
    """
    x = 0
    while True:
        yield x
        if x == state:
            break
        x = (x - state) & state


def genSubsetOfStateDescending(state: int):
    """降序枚举state所有子集的子集.
    0b1101 -> 13,12,9,8,5,4,1,0.
    """
    x = state
    while True:
        yield x
        if x == 0:
            break
        x = (x - 1) & state


def genSupersetOfState(n: int, state: int):
    """升序枚举state的所有超集.
    0b1101 -> 13,15.
    """
    upper = 1 << n
    x = state
    while x < upper:
        yield x
        x = (x + 1) | state


def genSubsetOfSizeK(n: int, k: int):
    """遍历n个元素的集合中大小为k的子集(combinations).
    一共有C(n,k)个子集.
    C(4,2) -> 3,5,6,9,10,12.
    """
    if k <= 0 or k > n:
        return
    upper = 1 << n
    x = (1 << k) - 1
    while x < upper:
        yield x
        t = x | (x - 1)
        # nextCombination (gosper hack)
        x = (t + 1) | (((~t & -~t) - 1) >> ((x & -x).bit_length()))


if __name__ == "__main__":
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

    for sub in genSubsetOfStateAscending(0b1101):
        print(bin(sub)[2:])
    for sub in genSubsetOfStateDescending(0b1101):
        print(bin(sub)[2:])
    for sup in genSupersetOfState(4, 0b1101):
        print(bin(sup)[2:])
    for sub in genSubsetOfSizeK(4, 2):
        print(bin(sub)[2:], "999")
