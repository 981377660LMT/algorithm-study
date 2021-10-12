/**
 * @param {string} IP
 * @return {string}
 * 编写一个函数来验证输入的字符串是否是有效的 IPv4 或 IPv6 地址。
 */
const validIPAddress = function (IP: string): string {
  const arr4 = IP.split('.')
  const arr6 = IP.split(':')
  if (arr4.length === 4) {
    // 如果用正则表达式判断每组数小于256比较繁杂
    // 这里先用正则判断是否为3位数字以内
    // IPv4不支持开头的0，比如 01.01.01.01 是无效的。
    // 再判断数字是否小于256即可
    if (arr4.every(part => part.match(/^0$|^([1-9]\d{0,2})$/) && Number(part) <= 255)) {
      return 'IPv4'
    }
  } else if (arr6.length === 8) {
    // 直接判断8组数均为4位以内16进制数即可
    // IPv6允许打头的0存在。
    if (arr6.every(part => part.match(/^[0-9a-fA-F]{1,4}$/))) {
      return 'IPv6'
    }
  }
  return 'Neither'
}

// 该题目要三种主要解法：

// 正则表达式，该方法性能不太好。
// 分治法，效率最高的方法之一。
// 使用分治法和内置的 try/catch，将字符串转换成整数处理。
// 使用 try/catch 不是一种好的方式，因为 try 块中的代码不会被编译器优化，所以最好不要在面试中使用。

// IPv4 地址由十进制数和点来表示，每个地址包含 4 个十进制数，其范围为 0 - 255， 用(".")分割。
// 同时，IPv4 地址内的数不会以 0 开头

// IPv6 地址由 8 组 16 进制的数字来表示，每组表示 16 比特
// 而且，我们可以加入一些以 0 开头的数字，字母可以使用大写，也可以是小写。
// 所以， 2001:db8:85a3:0:0:8A2E:0370:7334 也是一个有效的 IPv6 address地址 (即，忽略 0 开头，忽略大小写)。
