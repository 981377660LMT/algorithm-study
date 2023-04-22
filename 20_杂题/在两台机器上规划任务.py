# https://cp-algorithms.com/schedules/schedule_two_machines.html
# 在两台机器上分配任务/调度任务/规划任务
# 如果我们有两台以上的机器，这个任务就变成了 NP 完全问题。


# 给定n个任务和两台机器，每个任务需要在两台机器上处理，
# 第i个任务需要在第一台机器上花费ai时间，在第二台机器上花费bi时间。每台机器一次只能处理一个任务。求最优的调度方案。
# !对任务按照 min(ai,bi) 进行排序
# !min(aj,bj+1)<=min(bj,aj+1)


from typing import List, Tuple


def scheduleTwoMachines(jobs: List[Tuple[int, int]]) -> Tuple[int, int]:
    jobsWithId = [(a, b, i) for i, (a, b) in enumerate(jobs)]
    jobsWithId.sort(key=lambda x: min(x[0], x[1]))
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
        time2 = max(time1, time2) + job[1]
    return time1, time2


if __name__ == "__main__":
    jobs = [(3, 1), (1, 3), (2, 2), (2, 1), (1, 2)]
    print(scheduleTwoMachines(jobs))
