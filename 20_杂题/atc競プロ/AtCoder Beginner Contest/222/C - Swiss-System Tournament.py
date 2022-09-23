"""剪刀石头布巡回赛
Swiss-System Tournament(瑞士制锦标赛)
与小组赛制或其他从比赛一开始就知道所有配对的系统不同，在瑞士系统中，
每轮比赛的配对是在上一轮结束后完成的，并取决于其结果。
"""


RULE = {
    ("G", "G"): (0, 0),
    ("G", "C"): (1, 0),
    ("G", "P"): (0, 1),
    ("C", "C"): (0, 0),
    ("C", "G"): (0, 1),
    ("C", "P"): (1, 0),
    ("P", "P"): (0, 0),
    ("P", "G"): (1, 0),
    ("P", "C"): (0, 1),
}

if __name__ == "__main__":

    def pk(player1: int, type1: str, player2: int, type2: str) -> None:
        """锤子剪刀布 G:锤子 C:剪刀 P:布"""
        score1, score2 = RULE[(type1, type2)]  # type: ignore
        win[player1] += score1
        win[player2] += score2

    n, m = map(int, input().split())
    grid = [input() for _ in range(n * 2)]

    win = [0] * n * 2
    players = list(range(n * 2))
    for round in range(m):
        for p1, p2 in zip(players[::2], players[1::2]):  # 每轮根据当前实时排名进行pk
            pk(p1, grid[p1][round], p2, grid[p2][round])
        players = sorted(range(2 * n), key=lambda x: (-win[x], x))

    print(*[num + 1 for num in players], sep="\n")
