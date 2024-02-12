# https://codeforces.com/problemset/problem/896/C
# ->
# https://www.luogu.com.cn/problem/CF896C


import os


def getFilePath(name: str) -> str:
    return os.path.join(os.path.dirname(__file__), name)


inputPath = getFilePath("a.txt")
outputPath = getFilePath("b.txt")
prefix = "https://www.luogu.com.cn/problem/CF"

with open(inputPath, "r") as reader:
    with open(outputPath, "a") as writer:
        for line in reader:
            line = line.strip()
            words = line.split("/")
            if len(words) > 2:
                newLine = f"{prefix}{words[-2]}{words[-1]}"
                writer.write(newLine + "\n")
