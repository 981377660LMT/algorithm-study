# SelectPow2


from typing import List


def selectPow2(k: int, target: int) -> List[int]:
    """
    从1,2,4,...中选出k个数,每个数可以选择无数次(可以组成[k,k*(2^n)]中任意一个数).
    返回一组选择方案,使得和为target,如果不存在,返回空列表.
    """
    if k > target or k < target.bit_count():
        return []

    counter = [0] * 100
    count = 0
    for bit in range(len(counter)):
        if target & (1 << bit):
            counter[bit] = 1
            count += 1

    for i in range(len(counter) - 1, -1, -1):
        if i > 0:
            diff = k - count
            min_ = min(diff, counter[i])
            counter[i] -= min_
            counter[i - 1] += min_ * 2
            count += min_

    res = []
    for i in range(len(counter)):
        res += [1 << i] * counter[i]
    return res


if __name__ == "__main__":
    print(selectPow2(4, 8))
