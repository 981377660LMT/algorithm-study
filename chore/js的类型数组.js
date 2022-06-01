// ArrayBuffer对象保存着原始二进制数据，它不能对数据进行操作，
// 只能通过视图类（TypeArray或DataView）才能对数据进行读写。
const buffer = new ArrayBuffer(4 * 1000)
const hashTable1 = new Int32Array(buffer)
// 两个效果一样
const hashTable2 = new Int32Array(1000)

// Float32Array 等价于C里的float
// Float64Array 等价于C里的double
// 但注意js里的number默认是64位的double

// If control over byte order is needed, use DataView instead.
const view = new DataView(new ArrayBuffer(4))
view.setUint32(0, 0x12345678) // 默认big-endian(人的认知)
console.log(view)

// IEEE-754 标准的 64 bit 双精度浮点数(double):
// sign(符号): 占 1 bit, 表示正负;
// exponent(指数): 占 11 bit，表示范围;
// mantissa(尾数): 占 52 bit，表示精度，多出的末尾如果是 1 需要进位;

// Java的Double为什么会丢失精度
// 谈到有小数点的加减乘除都会想到用Java里的BigDecimal来解决
// float：2^23 = 8388608，一共七位，这意味着最多能有7位有效数字，但绝对能保证的为6位，也即float的精度为6~7位有效数字；
// double：2^52 = 4503599627370496，一共16位，同理，double的精度为15~16位。

console.log((1 + 2) / 10)
console.log(0.1 + 0.2) // 0.30000000000000004 python里解决方法为Fraction
// 解决方法还可以用字符串运算

// >>> # Don't do this:
// >>> 0.1 + 0.2 == 0.3
// False

// >>> # Do this instead:
// >>> math.isclose(0.1 + 0.2, 0.3)
// True
