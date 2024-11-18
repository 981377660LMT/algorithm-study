import sys

sys.setrecursionlimit(1 << 25)

N, M, L = map(int, input().split())
A = list(map(int, input().split()))
B = list(map(int, input().split()))
C = list(map(int, input().split()))
cards = A + B + C
number_of_cards = len(cards)

t_hand = 0
for i in range(N):
    t_hand |= 1 << i

a_hand = 0
for i in range(N, N + M):
    a_hand |= 1 << i

table = 0
for i in range(N + M, N + M + L):
    table |= 1 << i

memo = {}


def can_win(t_hand, a_hand, table, turn):
    key = (t_hand, a_hand, table, turn)
    if key in memo:
        return memo[key]

    if turn == 0:
        if t_hand == 0:
            memo[key] = False
            return False
        for i in range(number_of_cards):
            if (t_hand >> i) & 1:
                new_t_hand = t_hand & ~(1 << i)
                new_table = table | (1 << i)
                pick_options = [
                    j
                    for j in range(number_of_cards)
                    if (new_table >> j) & 1 and cards[j] < cards[i]
                ]
                pick_options.append(None)
                for pick in pick_options:  # do not pick
                    temp_t_hand = new_t_hand
                    temp_table = new_table
                    if pick is not None:
                        temp_t_hand |= 1 << pick
                        temp_table &= ~(1 << pick)
                    if not can_win(temp_t_hand, a_hand, temp_table, 1):
                        memo[key] = True
                        return True
        memo[key] = False
        return False
    else:
        if a_hand == 0:
            memo[key] = False
            return False
        for i in range(number_of_cards):
            if (a_hand >> i) & 1:
                new_a_hand = a_hand & ~(1 << i)
                new_table = table | (1 << i)
                pick_options = [
                    j
                    for j in range(number_of_cards)
                    if (new_table >> j) & 1 and cards[j] < cards[i]
                ]
                pick_options.append(None)
                for pick in pick_options:
                    temp_a_hand = new_a_hand
                    temp_table = new_table
                    if pick is not None:
                        temp_a_hand |= 1 << pick
                        temp_table &= ~(1 << pick)
                    if not can_win(t_hand, temp_a_hand, temp_table, 0):
                        memo[key] = True
                        return True
        memo[key] = False
        return False


if can_win(t_hand, a_hand, table, 0):
    print("Takahashi")
else:
    print("Aoki")
