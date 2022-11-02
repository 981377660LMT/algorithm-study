# F - Hammer 2
# 锤子与墙壁

# 从原点开始，到目标点target
# 求最少的路程。
# 有n个锤子，n个墙
# w[i] 第i面墙的位置
# h[i] 第i个锤子的位置
# 第i面墙只能被第i个锤子打碎

# !n<=1500
# !纯模拟解法
# 1.先把墙壁和锤子的位置处理出来
# 2.然后向左最远 向右最远 向左最远 向右最远
# 3.记录当前最左走到哪里 最右走到哪里 当前的位置 下一步需要到达的位置

INF = int(1e18)


def main():
    _, goal = map(int, input().split())  # 目的地
    walls = list(map(int, input().split()))
    hammers = list(map(int, input().split()))

    rightEvents = []  # 原点右侧的(墙, 锤子)对
    leftEvents = []  # 原点左侧的(墙, 锤子)对
    for wall, hammer in zip(walls, hammers):
        if (wall < hammer < 0) or (0 < hammer < wall):  # !剪枝 没用的墙壁
            # ハンマーが壁の手前にある、無意味な壁はハンマーとともに除去
            continue
        if (hammer < wall < 0) or (0 < wall < hammer):  # !无法打碎的墙壁
            # ハンマーが壁の奥にある、壊せない壁のハンマーだけ除去
            hammer = 0
        rightEvents.append((wall, hammer))
    rightEvents = sorted(rightEvents, key=lambda x: x[0])
    # 壁の位置で分割
    for i in range(len(rightEvents)):
        if rightEvents[i][0] > 0:
            leftEvents = rightEvents[:i][::-1]
            del rightEvents[:i]
            break
    else:
        leftEvents = rightEvents[::-1]
        del rightEvents[:]

    # res = 必要な移動距離
    # goal_min = 走査限界位置の最小値
    # goal_max = 走査限界位置の最大値
    res = abs(goal)

    goal_min = -INF
    goal_max = INF

    while True:
        # goal_next = 次に必要なハンマーの位置
        goal_next = 0  # !下一个必要的锤子的位置

        if goal < 0:
            # 目的位置がマイナス

            # 原点から近い順で、目標の位置まで処理
            for wall, hammer in leftEvents:
                if wall < goal:
                    break

                # 壊せない壁に当たったら'-1'で終了
                if hammer == 0:
                    print(-1)
                    return

                # より遠くのハンマーとその壁を記録
                if hammer > goal_next:
                    goal_next = hammer
                    goal_min = wall

            # 走査限界位置より遠い位置に必要なハンマーがあった場合は'-1'で終了
            if goal_next > goal_max:
                print(-1)
                return

        else:
            # 目的位置がプラス（以下マイナスと同じ処理）

            for wall, hammer in rightEvents:
                if wall > goal:
                    break
                if hammer == 0:
                    print(-1)
                    return
                if hammer < goal_next:
                    goal_next = hammer
                    goal_max = wall
            if goal_next < goal_min:
                print(-1)
                return

        # 次に必要なハンマーが無かった場合、今までに必要な移動距離を出力して終了
        if goal_next == 0:
            print(res)
            return

        # 必要な移動距離を往復分で足す
        res += 2 * abs(goal_next)
        # 次に必要なハンマーの位置を次の目標にする
        goal = goal_next


main()
