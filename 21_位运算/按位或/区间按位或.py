# 区间按位或
# https://www.geeksforgeeks.org/bitwise-or-or-of-a-range/


# 1. Find the position of Most Significant Bit (MSB) in both the numbers (L and R)
# 2. If the position of both MSBs are different, set all the bits from the max(MSB1, MSB2),
#    including this different bit, up to Oth bit i.e., add the value (1 << i) for all 0 ≤ i ≤ max(MSB1, MSB2) in the answer.
# 3. If the position of both MSBs is the same, then
#     - Set this bit corresponding to MSB or add the value (1 << MSB) in the answer.
#     - Subtract the value (1 << MSB) from both the numbers (L and R).
#     - Repeat steps 1, 2, and 3.


def rangeOr(left: int, right: int) -> int:
    assert left <= right
    msb1, msb2 = left.bit_length(), right.bit_length()
    if msb1 != msb2:
        return (1 << msb2) - 1

    res = 0
    while msb1 == msb2:
        cand = 1 << msb1
        left ^= cand
        right ^= cand
        res |= cand
        msb1, msb2 = left.bit_length(), right.bit_length()

    max_ = max(msb1, msb2)
    return res | ((1 << max_) - 1)


if __name__ == "__main__":

    # check with brute force:

    import random
    from functools import reduce

    while True:
        left = random.randint(0, 10000)
        right = random.randint(left, 10000)
        if rangeOr(left, right) != reduce(lambda x, y: x | y, range(left, right + 1), 0):
            print(left, right)
            break
