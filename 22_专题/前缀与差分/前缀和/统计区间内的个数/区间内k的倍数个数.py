# 区间内k的倍数个数
def numberOfMultiples(lower: int, upper: int, k: int) -> int:
    return upper // k - (lower - 1) // k


print(numberOfMultiples(1, 10, 2))
