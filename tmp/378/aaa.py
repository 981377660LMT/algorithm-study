import time


def genRandomString(length):
    import random

    return "".join(random.choice("abcdefghijklmnopqrstuvwxyz") for _ in range(length))


ss = [genRandomString(2000) for _ in range(2000)]

time1 = time.time()
for s in ss:
    b = s[len(s) // 2 :][::-1]
    # b = s[len(s) - 1 : len(s) // 2 - 1 : -1]

print(time.time() - time1)
