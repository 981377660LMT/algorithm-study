# 给定n，问有多少正整数组(A,B,C,D)，满足 AB+CD=n
# n<=2e5

# nlogn预处理A*B=i的方案数即可


def solve(n: int) -> int:
    counter = [0] * (n + 1)  # counter[i] = A*B=i的方案数
    for a in range(1, n + 1):
        for b in range(1, n // a + 1):
            counter[a * b] += 1
    return sum(counter[a] * counter[n - a] for a in range(1, n))


if __name__ == "__main__":
    n = int(input())
    print(solve(n))
