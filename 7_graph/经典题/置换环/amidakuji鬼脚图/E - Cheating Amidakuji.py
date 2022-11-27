"""
amidakuji 鬼脚图

为啥叫 cheating amidakuji
按规则每次走到横线就要移动过去 (横线对应题目里的数组)
有一次破坏横线的机会 所以叫cheating
(1表示的是当たり的签)
"""

from typing import List


def cheatingAmidakuji(n: int, lines: List[int]) -> List[int]:
    """
    Args:
        n: 鬼脚图的n个出发点1-n
        lines: 鬼脚图的横线,一共m条,从上往下表示.
               lines[i] 表示第i条横线连接 line[i] 和 line[i]+1.
    Returns:
        有m个人抽签作弊,第i个人不执行第i条横线的交换(移除第i条横线).
        问每个人的签中1在终点的位置(1就是あみだくじ里的`当たり`)
    """
    perm = list(range(1, n + 1))
    for col in lines:
        perm[col - 1], perm[col] = perm[col], perm[col - 1]
    # !执行`完整交换`后 每个元素最后到了哪个位置(加速模拟) 0-index
    mp = {num: i for i, num in enumerate(perm)}

    perm = list(range(1, n + 1))
    res = []
    for col in lines:
        if perm[col - 1] == 1:  # 对结果有影响
            res.append(mp[perm[col]])
        elif perm[col] == 1:  # 对结果有影响
            res.append(mp[perm[col - 1]])
        else:
            res.append(mp[1])  # 对结果无影响
        perm[col - 1], perm[col] = perm[col], perm[col - 1]
    return res


if __name__ == "__main__":
    n, m = map(int, input().split())
    nums = list(map(int, input().split()))
    res = cheatingAmidakuji(n, nums)
    print(*[i + 1 for i in res], sep="\n")
