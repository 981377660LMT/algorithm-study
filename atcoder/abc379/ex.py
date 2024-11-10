import sys
import threading

MOD = 998244353


def main():
    import sys
    from itertools import product

    sys.setrecursionlimit(1 << 25)

    H, W = map(int, sys.stdin.readline().split())
    S = [sys.stdin.readline().strip() for _ in range(H)]

    if H < W:
        H, W = W, H
        S = list(map(list, zip(*S)))

    def get_options(row):
        options = []
        for c in row:
            if c == "?":
                options.append([0, 1, 2])
            else:
                options.append([int(c) - 1])
        return options

    def generate_valid_rows(options):
        valid = []
        state = []

        def back(pos, prev_color, current_state):
            if pos == W:
                valid.append(current_state)
                return
            for color in options[pos]:
                if color != prev_color:
                    back(pos + 1, color, current_state * 3 + color)

        back(0, -1, 0)
        return valid

    valid_rows = [generate_valid_rows(get_options(row)) for row in S]

    compat = []
    for i in range(H):
        compat.append([])
        for current in valid_rows[i]:
            if i == 0:
                compat[i].append(None)
            else:
                compat[i].append([])
                for prev in valid_rows[i - 1]:
                    conflict = False
                    tmp_curr = current
                    tmp_prev = prev
                    for _ in range(W):
                        if tmp_curr % 3 == tmp_prev % 3:
                            conflict = True
                            break
                        tmp_curr = tmp_curr // 3
                        tmp_prev = tmp_prev // 3
                    if not conflict:
                        compat[i][-1].append(prev)
        if i > 0:
            for idx, current in enumerate(valid_rows[i]):
                compat[i][idx] = []
                for prev in valid_rows[i - 1]:
                    tmp_curr = current
                    tmp_prev = prev
                    conflict = False
                    for _ in range(W):
                        if tmp_curr % 3 == tmp_prev % 3:
                            conflict = True
                            break
                        tmp_curr = tmp_curr // 3
                        tmp_prev = tmp_prev // 3
                    if not conflict:
                        compat[i][idx].append(prev)

    dp_prev = {}
    for row in valid_rows[0]:
        dp_prev[row] = 1

    for i in range(1, H):
        dp_current = {}
        for current in valid_rows[i]:
            total = 0
            for prev in compat[i][valid_rows[i].index(current)]:
                total = (total + dp_prev.get(prev, 0)) % MOD
            dp_current[current] = total
        dp_prev = dp_current

    result = sum(dp_prev.values()) % MOD if H > 0 else 0
    print(result)


main()
