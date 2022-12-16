# 给定初始位置X,给定每步向左或向右可移动距离D,
# 给定移动的步数K,问执行完这些步数后，
# 新位置的绝对值最小是多少，输出这个最小值。

# !分类讨论
def walkingTakahashi(start: int, times: int, step: int) -> int:
    if start < 0:
        start = -start
    if start >= times * step:
        return start - times * step  # 无法移动到0,只能移动到最近的位置
    div = start // step
    pos = start - div * step
    remain = times - div

    if remain & 1:  # 还需要移动奇数次
        return step - pos
    return pos


if __name__ == "__main__":
    X, K, D = map(int, input().split())
    print(walkingTakahashi(X, K, D))
