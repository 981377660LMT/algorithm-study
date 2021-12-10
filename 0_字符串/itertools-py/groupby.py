from itertools import groupby

meetings = [[1, 2, 5], [2, 3, 8], [1, 5, 10], [2, 6, 10]]


for key, group in groupby(meetings, key=lambda x: x[2]):
    print(key, group)
    for x, y, _ in group:
        print(x, y, _)

