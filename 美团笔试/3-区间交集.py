# 为了使多台电脑上文件的版本一致，小团使用了一个版本管理软件进行同步。
# 该软件会比较文件之间的差别并汇报冲突，
# 所有冲突都被解决才能完成同步。比如在上一次同步之后，
# 小团将A机器上的文件f修改成了版本f1，
# 并在f1没有同步到B机器上时将 B机器上的文件f修改成了版本f2，
# 则版本管理软件会检测到这一冲突并汇报给小团。
# 否则版本管理软件会直接使用新版的文件覆盖另一台机器上的旧版文件，
# 此时不会汇报冲突。

# 在上一次同步之后，小团修改了他的游戏本上的许多文件和老爷机上的许多文件。
# 现在给出这些被修改的文件的编号，请你求出版本管理软件会汇报多少个文件的冲突。


# 。文件编号为1到n n<=1e10
while True:
    try:
        n, m1, m2 = map(int, input().split())
        starts1 = list(map(int, input().split()))
        ends1 = list(map(int, input().split()))
        starts2 = list(map(int, input().split()))
        ends2 = list(map(int, input().split()))

        inter1 = []
        inter2 = []
        for s1, e1 in zip(starts1, ends1):
            inter1.append((s1, e1))
        for s2, e2 in zip(starts2, ends2):
            inter2.append((s2, e2))

        inter1.sort()
        inter2.sort()
        inter1.append((int(1e20), int(1e20)))
        res = 0
        index1 = 0

        for index2, (start2, end2) in enumerate(inter2):
            # 包含，相交，相离
            while index1 < len(inter1) + 1:
                start1, end1 = inter1[index1]
                # 线段相交长度
                left = max(start1, start2)
                right = min(end1, end2)
                if right - left < 0:
                    if end2 < start1:
                        break

                res += max(0, right - left + 1)
                index1 += 1

        print(res)

    except EOFError:
        break
