from math import floor
import sys


def get_tile(x, y, K):
    tile_i = (2 * x + 1) // (2 * K)
    tile_j = (2 * y + 1) // (2 * K)

    if (tile_i % 2) == (tile_j % 2):
        k = y - tile_j * K
    else:
        k = x - tile_i * K

    return (tile_i, tile_j, k)


def main():
    input = sys.stdin.readline
    T = int(input())
    results = []

    for _ in range(T):

        K, Sx, Sy, Tx, Ty = map(int, input().split())
        i_s, j_s, k_s = get_tile(Sx, Sy, K)
        i_t, j_t, k_t = get_tile(Tx, Ty, K)

        parity_s = (i_s % 2) == (j_s % 2)
        parity_t = (i_t % 2) == (j_t % 2)

        if i_s == i_t and j_s == j_t:
            steps = abs(k_s - k_t)
        else:
            delta_i = abs(i_s - i_t)
            delta_j = abs(j_s - j_t)
            steps = delta_i + delta_j
            if (parity_s != parity_t) and (k_s != k_t):
                steps += 1
        results.append(str(steps))

    print("\n".join(results))


if __name__ == "__main__":
    main()
