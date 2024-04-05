# https://cp-algorithms.com/schedules/schedule_two_machines.html
# !机器A串行，机器B串行
# 在两台机器上分配任务/调度任务/规划任务
# 如果我们有两台以上的机器，这个任务就变成了 NP 完全问题。
# 流水调度问题(JOHNSON贪心算法)
# 给定n个任务和两台机器，每个任务需要在两台机器上处理，
# 第i个任务需要在第一台机器上花费ai时间，在第二台机器上花费bi时间。每台机器一次只能处理一个任务。求最优的调度方案。
# !对任务按照 min(ai,bi) 进行排序
# !min(aj,bj+1)<=min(bj,aj+1)


from typing import List, Tuple


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


def scheduleTwoMachines(jobs: List[Tuple[int, int]]) -> Tuple[int, int, List[int]]:
    """
    :param jobs: [(a1, b1), (a2, b2), ...] 任务列表
    :return: (time1, time2, order) 机器1的总时间,机器2的总时间,任务的顺序
    """
    jobsWithId = [(a, b, i) for i, (a, b) in enumerate(jobs)]
    jobsWithId.sort(key=lambda x: min2(x[0], x[1]))
    jobs1, jobs2 = [], []
    for job in jobsWithId:
        a, b = job[0], job[1]
        if a < b:
            jobs1.append(job)
        else:
            jobs2.append(job)
    newJobs = jobs1 + jobs2[::-1]

    time1, time2 = 0, 0
    for job in newJobs:
        time1 += job[0]
        time2 = max2(time1, time2) + job[1]
    return time1, time2, [job[2] for job in newJobs]


if __name__ == "__main__":
    n = int(input())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    time1, time2, order = scheduleTwoMachines(list(zip(A, B)))
    print(time2)
    print(" ".join(str(x + 1) for x in order))
