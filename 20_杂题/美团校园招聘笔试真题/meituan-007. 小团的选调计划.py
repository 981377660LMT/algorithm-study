n = int(input())

wants = []
for _ in range(n):
    line = [int(v) - 1 for v in input().split()]
    wants.append(line)

visited = [False] * n
for wa in wants:
    for cand in wa:
        if not visited[cand]:
            visited[cand] = True
            print(cand + 1, end=' ')
            break

# 注意输出 end 分隔符
