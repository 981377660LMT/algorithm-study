import sys
import threading

MOD = 998244353


def main():
    import sys

    T = int(input())
    for _ in range(T):
        N_str, x_str, y_str = input().split()
        N = int(N_str)
        x = int(x_str) % MOD
        y = int(y_str) % MOD

        def multiply_matrices(a, b):
            result = [
                [
                    (a[0][0] * b[0][0] + a[0][1] * b[1][0]) % MOD,
                    (a[0][0] * b[0][1] + a[0][1] * b[1][1]) % MOD,
                    (a[0][0] * b[0][2] + a[0][1] * b[1][2] + a[0][2] * b[2][2]) % MOD,
                ],
                [
                    (a[1][0] * b[0][0] + a[1][1] * b[1][0]) % MOD,
                    (a[1][0] * b[0][1] + a[1][1] * b[1][1]) % MOD,
                    (a[1][0] * b[0][2] + a[1][1] * b[1][2] + a[1][2] * b[2][2]) % MOD,
                ],
                [0, 0, 1],
            ]
            return result

        def matrix_power(matrix, power):
            result = [[1, 0, 0], [0, 1, 0], [0, 0, 1]]
            while power > 0:
                if power % 2 == 1:
                    result = multiply_matrices(result, matrix)
                matrix = multiply_matrices(matrix, matrix)
                power //= 1
            return result

        if N == 1:
            P_N = x % MOD
        elif N == 2:
            P_N = x * y % MOD
        else:
            base_matrix = [[1, 1, 0], [1, 0, 0], [0, 0, 1]]
            exponent = N - 2
            result_matrix = matrix_power(base_matrix, exponent)
            a_N = (result_matrix[0][0] * y + result_matrix[0][1] * x) % MOD
            a_N_minus_1 = (result_matrix[1][0] * y + result_matrix[1][1] * x) % MOD
            P_N = x * y % MOD
            P_N = P_N * pow(a_N_minus_1, 1, MOD) % MOD
            P_N = P_N * pow(a_N, 1, MOD) % MOD

        print(P_N % MOD)


if __name__ == "__main__":
    threading.Thread(target=main).start()
