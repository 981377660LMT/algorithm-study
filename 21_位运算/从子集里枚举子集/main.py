from typing import List


def main(nums: List[int]):
    n = len(nums)
    res = []
    for state in range(1 << n):
        group1, group2 = state, 0
        while group1:
            res.append((group1, group2))
            # 关键，不断减一+与运算跳数
            group1 = state & (group1 - 1)
            group2 = state ^ group1
    return res


if __name__ == '__main__':
    print(main([1, 2, 3, 4]))

