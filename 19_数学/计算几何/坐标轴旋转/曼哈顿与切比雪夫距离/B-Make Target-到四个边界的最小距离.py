# B - Make Target(同心矩形)
# https://atcoder.jp/contests/abc395/tasks/abc395_b

# n=2
##
##

# n=5
#####
# ...#
# .#.#
# ...#
#####

# n=8
########
# ......#
# .####.#
# .#..#.#
# .#..#.#
# .####.#
# ......#
########


from typing import List


def makeTarget(n: int) -> List[str]:
    """生成n阶靶型图案(target)."""
    target = [["?"] * n for _ in range(n)]
    for i in range(n):
        for j in range(n):
            minDistToBorder = min(i, j, n - i - 1, n - j - 1)
            target[i][j] = "." if minDistToBorder & 1 else "#"
    return ["".join(row) for row in target]


if __name__ == "__main__":
    n = int(input())
    target = makeTarget(n)
    print("\n".join(target))
