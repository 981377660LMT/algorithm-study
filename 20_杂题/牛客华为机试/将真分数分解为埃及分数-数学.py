# 分子为1的分数称为埃及分数。现输入一个真分数(分子比分母小的分数，叫做真分数)，请将该分数分解为埃及分数。
# 如：8/11 = 1/2+1/5+1/55+1/110。
while True:
    try:
        arr = input().split('/')
        fz = int(arr[0])
        c = []
        for _ in range(fz):
            c.append(f'{1}' + '/' + arr[1])
        print('+'.join(c))

    except EOFError:
        break
