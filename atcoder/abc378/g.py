def main():
    import sys
    import threading

    def solve():

        A, B, M = map(int, sys.stdin.readline().split())
        N = A * B - 1

        # Precompute factorials modulo M
        factorial = [1] * (N + 1)
        for i in range(1, N + 1):
            factorial[i] = factorial[i - 1] * i % M

        inv_factorial = [1] * (N + 1)
        inv_factorial[N] = pow(factorial[N], M - 2, M)
        for i in range(N - 1, -1, -1):
            inv_factorial[i] = inv_factorial[i + 1] * (i + 1) % M

        total = 0

        def generate_partitions_with_first_row(N, A):
            def helper(remaining, current_partition):
                if remaining == 0:
                    yield current_partition
                    return
                start = min(current_partition[-1], remaining)
                for k in range(start, 0, -1):
                    if k <= A:
                        yield from helper(remaining - k, current_partition + [k])

            yield from helper(N - A, [A])

        partitions = list(generate_partitions_with_first_row(N, A))

        for partition in partitions:
            # Compute conjugate partition
            max_col = max(partition)
            λ_prime = [0] * max_col
            for length in partition:
                for j in range(length):
                    λ_prime[j] += 1
            # Check if λ'[0] == B
            if λ_prime[0] != B:
                continue
            # Compute hook lengths
            hook_lengths = []
            for i in range(len(partition)):
                for j in range(partition[i]):
                    h = partition[i] - j
                    k = 0
                    for l in range(i + 1, len(partition)):
                        if partition[l] > j:
                            k += 1
                        else:
                            break
                    h += k
                    hook_lengths.append(h)
            # Compute f^λ modulo M
            numerator = factorial[N]
            denominator = 1
            for h in hook_lengths:
                denominator = denominator * h % M
            f_lambda = numerator * pow(denominator, M - 2, M) % M
            total = (total + f_lambda * f_lambda) % M

        print(total)

    threading.Thread(target=solve).start()


if __name__ == "__main__":
    main()
