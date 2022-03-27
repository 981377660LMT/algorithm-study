while True:
    try:
        m, n = map(int, input().split())
        res = []
        for num in range(m, n + 1):
            chars = [int(char) for char in str(num)]
            sum_ = sum(c ** 3 for c in chars)
            if sum_ == num:
                res.append(num)

        if not res:
            print('no')
        else:
            for num in res:
                print(num, end=' ')
    except EOFError:
        break

