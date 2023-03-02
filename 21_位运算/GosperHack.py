# GosperHack 枚举大小k的子集


def enumerationCombinations(n, k):
    """枚举大小为k的子集,其中n为集合的大小."""
    s = (1 << k) - 1
    while s < (1 << n):  # nextCombination
        yield s
        x = s & -s
        y = s + x
        s = ((s & ~y) // x >> 1) | y


for i in enumerationCombinations(5, 3):
    print(bin(i))
