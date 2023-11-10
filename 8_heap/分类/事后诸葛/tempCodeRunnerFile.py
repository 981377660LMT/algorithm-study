
                for q in range(2):
                    points = [randint(-10, 10) for _ in range(n)]
                    initHps = [randint(0, 10) for _ in range(q)]
                    res1 = solve(points, initHps)
                    res2 = bf(points, initHps)
                    if res1 != res2:
                        print(points, initHps)
                        print(res1, res2)
                        return