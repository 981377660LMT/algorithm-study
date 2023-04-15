# https://kopricky.github.io/code/Academic/gale_sharpley.html
# 稳定婚配问题(结婚问题)
# O(n^2)


from typing import List


def galeSharpley(boys: List[List[int]], girls: List[List[int]]) -> List[int]:
    """
    给出一个n个男生和n个女生的列表,其中:
    每位男生根据对所有女生的心仪程度从高到低进行排名,
    每位女生根据对所有男生的心仪程度从高到低进行排名.
    找到一个"稳定匹配".

    返回每位男生的配对女生的编号.
    """
    n = len(boys)
    worder = [[0] * n for _ in range(n)]
    single = []
    for i in range(n):
        single.append(i)
        for j in range(n):
            worder[i][girls[i][j]] = j
    idPos = [0] * n
    gToB = [-1] * n
    while single:
        boy = single.pop()
        while True:
            girl = boys[boy][idPos[boy]]
            idPos[boy] += 1
            if gToB[girl] < 0:
                gToB[girl] = boy
                break
            elif worder[girl][gToB[girl]] > worder[girl][boy]:
                single.append(gToB[girl])
                gToB[girl] = boy
                break

    bToG = [0] * n
    for i in range(n):
        bToG[gToB[i]] = i
    return bToG


if __name__ == "__main__":
    # 索引代表男性编号（0，1，2），数组代表喜好女性的列表
    boys = [[0, 1, 2], [1, 0, 2], [0, 1, 2]]
    # 索引代表女性编号（0，1，2），数组代表喜好男性的列表
    girls = [[1, 2, 0], [0, 1, 2], [0, 1, 2]]
    print(galeSharpley(boys, girls))  # [1, 0, 2]
