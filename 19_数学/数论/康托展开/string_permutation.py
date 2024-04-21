# string_permutation


# s を並び替えてできる文字列のうち，s が辞書順で何番目か(0-indexed)を求める
def find_string_no_permutation_no(s):
    n = len(s)

    fact = [1] * (n + 1)
    for i in range(1, n + 1):
        fact[i] = fact[i - 1] * i  # !mod

    freq = dict()
    for c in s:
        if c not in freq:
            freq[c] = 0
        freq[c] += 1
    chars = list(sorted(list(set(s))))

    res = 0
    for i in range(n):
        for c in chars:
            if freq[c] == 0:
                continue

            if c >= s[i]:
                break

            freq[c] -= 1

            num = fact[n - i - 1]
            for v in freq.values():
                if v != 0:
                    num //= fact[v]  # !mod inv

            res += num
            freq[c] += 1

        freq[s[i]] -= 1

    return res


# s を並び替えてできる文字列のうち，辞書順で k 番目(0-indexed)の文字列を求める
# 存在しない場合は空文字列を返す
def find_kth_string_permutation(k, s):
    n = len(s)

    fact = [1] * (n + 1)
    for i in range(1, n + 1):
        fact[i] = fact[i - 1] * i

    freq = dict()
    for c in s:
        if c not in freq:
            freq[c] = 0
        freq[c] += 1
    chars = list(sorted(list(set(s))))

    # k が作れるかチェック
    max_k = fact[n]
    for v in freq.values():
        max_k //= fact[v]

    if k >= max_k:
        return []

    res = []
    for i in range(n):
        # i 文字目を c にできるか
        for c in chars:
            if freq[c] <= 0:
                continue

            freq[c] -= 1

            # i 文字目を c にしたとき，作成できる文字列の個数
            count = fact[n - i - 1]
            for v in freq.values():
                count //= fact[v]

            if k < count:
                res.append(c)
                break

            freq[c] += 1
            k -= count

    return res


print(find_string_no_permutation_no(s="a" * int(1e3)))
print(find_kth_string_permutation(k=0, s="a" * int(1e3)))
