def countInRange(left: int, right: int, k: int) -> int:
    """
    左侧可选长度为L，右侧可选长度为R，长度不超过k的非空子数组个数.
    左侧和右侧都包含当前元素.
    """
    upper = right
    if upper > k:
        upper = k
    if upper <= 0:
        return 0
    if left > k:
        return (k + k - upper + 1) * upper // 2
    pos = k - left + 1
    if pos > upper:
        return upper * left
    c1 = pos - 1
    res1 = c1 * left
    c2 = upper - pos + 1
    min_ = k - (upper - 1)
    max_ = k - (pos - 1)
    res2 = (min_ + max_) * c2 // 2
    return res1 + res2


if __name__ == "__main__":

    def bruteForce(left: int, right: int, k: int) -> int:
        res = 0
        for L in range(1, left + 1):
            for R in range(1, right + 1):
                if L + R - 1 <= k:
                    res += 1
        return res

    import random

    for _ in range(10000):
        left = random.randint(1, 10)
        right = random.randint(1, 10)
        k = random.randint(1, 10)
        assert countInRange(left, right, k) == bruteForce(left, right, k)
    print("PASSED")
