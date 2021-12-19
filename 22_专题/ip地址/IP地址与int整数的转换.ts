// 将ip地址转换成10进制整数
function IPToInt(ip: string): number {
  const [s1, s2, s3, s4] = ip.split('.')
  const n1 = Number(s1) << 24
  const n2 = Number(s2) << 16
  const n3 = Number(s3) << 8
  const n4 = Number(s4) << 0
  return n1 | n2 | n3 | n4
}

function intToIP(int: number): string {
  const mask = 0xff
  const res: number[] = []
  res.push(int & mask)
  res.push((int >> 8) & mask)
  res.push((int >> 16) & mask)
  res.push((int >> 24) & mask)

  return res.join('')
}

console.log(IPToInt('10.0.3.193'))
console.log(intToIP(0b00001010_00000000_00000011_11000001))
