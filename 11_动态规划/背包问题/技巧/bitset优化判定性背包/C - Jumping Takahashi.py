from typing import List, Tuple


def jumpingTakahashi(jumps: List[Tuple[int, int]], target: int) -> bool:
    """
    起始时在原点.问跳跃n次是否能跳到target.
    每次跳跃jumps[i]可以跳jumps[i][0]或者跳jumps[i][1].
    n<=100, target<=1e4, jumps[i][0]<=100, jumps[i][1]<=100.
    时间复杂度 O(n*sum(jumps)/w)
    """
    mask = ~(-1 << (target + 1))
    dp = 1
    for a, b in jumps:
        dp = (dp << a) | (dp << b)
        dp &= mask
    return (dp >> target) & 1 != 0


if __name__ == "__main__":
    n, target = map(int, input().split())
    jumps = []
    for _ in range(n):
        jumps.append(tuple(map(int, input().split())))
    print("Yes" if jumpingTakahashi(jumps, target) else "No")
