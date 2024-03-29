// 墙上挂有 n 只已经打开的灯泡和 4 个按钮
// 假设这 n 只灯泡被编号为 [1, 2, 3 ..., n]，这 4 个按钮的功能如下：

// 将所有灯泡的状态反转（即开变为关，关变为开）
// 将编号为偶数的灯泡的状态反转
// 将编号为奇数的灯泡的状态反转
// 将编号为 3k+1 的灯泡的状态反转（k = 0, 1, 2, ...)

// 在进行了 presses 次未知操作后，你需要返回这 n 只灯泡可能有多少种不同的状态
// 前6个灯唯一地决定了其余的灯。这是因为每一个修改 第 x 的灯光的操作都会修改第(x+6) 的灯光。
// 前三个灯泡可以决定所有灯泡的状态：灯泡4的状态可以由1，2，3的状态线性表示
function flipLights(n: number, presses: number): number {}

// 参考
// !21_位运算/按位异或/线性基/672. 灯泡开关 Ⅱ.py
