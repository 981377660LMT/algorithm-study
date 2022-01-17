# 用户名的首字符必须是大写或者小写字母。
# 用户名只能包含大小写字母，数字。
# 用户名需要包含至少一个字母和一个数字。
# 如果用户名合法，请输出 "Accept"，反之输出 "Wrong"。

length = int(input())

for _ in range(length):
    isValid = True
    string = input()

    if not string[0].isalpha():
        isValid = False

    digitCount = 0
    alphaCount = 0

    if isValid:
        for char in string:
            if char.isalpha():
                alphaCount += 1
            elif char.isdigit():
                digitCount += 1
            else:
                isValid = False
                break

    if isValid and digitCount > 0 and alphaCount > 0:
        print("Accept")
    else:
        print("Wrong")

