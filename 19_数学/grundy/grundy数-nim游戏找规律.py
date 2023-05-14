# https://atcoder.jp/contests/abc297/tasks/abc297_g

# 状态太多导致无法用dfs求出初始状态的sg值
# !需要打表找规律,求出对应状态的sg值

# 初始时有n堆石头(nim游戏里的n堆石头)
# Alice和Bob轮流取石子，Alice先手
# 每个人必须取[L,R]个石子,不能取的人输
# 问Alice是否必胜


def nim(stone: int, left: int, right: int):
    """大小为stone的石头堆,每次取石子的数量范围为[left,right],对应的sg值."""
    return (stone % (left + right)) // left


if __name__ == "__main__":
    n, L, R = map(int, input().split())
    stones = list(map(int, input().split()))

    sg = 0
    for stone in stones:
        sg ^= nim(stone, L, R)
    print("First" if sg else "Second")
