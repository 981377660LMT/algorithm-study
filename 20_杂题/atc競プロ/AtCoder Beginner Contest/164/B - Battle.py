# 宝可梦战斗(对决)
# !注意这里可以二分回合数加速模拟(两只宝可梦的体力都大于0时的最大回合数)


if __name__ == "__main__":
    # 高橋君と青木君がモンスターを闘わせます。
    # 高橋君のモンスターは体力が A で攻撃力が B です。 青木君のモンスターは体力が C で攻撃力が D です。
    # 高橋君→青木君→高橋君→青木君→... の順に攻撃を行います。 攻撃とは、相手のモンスターの体力の値を自分のモンスターの攻撃力のぶんだけ減らすことをいいます。 このことをどちらかのモンスターの体力が 0 以下になるまで続けたとき、 先に自分のモンスターの体力が 0 以下になった方の負け、そうでない方の勝ちです。
    # 高橋君が勝つなら Yes、負けるなら No を出力してください。

    hp1, attack1, hp2, attack2 = map(int, input().split())

    # !二分回合数加速
    def check(mid: int) -> bool:
        return hp1 - attack2 * mid > 0 and hp2 - attack1 * mid > 0

    left, right = 1, int(1e16)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1

    # right 回合数内两只宝可梦的体力都大于0
    hp1 -= attack2 * right
    hp2 -= attack1 * right

    # !模拟
    while True:
        hp2 -= attack1
        if hp2 <= 0:
            print("Yes")
            break
        hp1 -= attack2
        if hp1 <= 0:
            print("No")
            break
