console.log('asas'.split('k'))

// 如果分隔符不存在 则返回原字符串
function numUniqueEmails(emails: string[]): number {
  function parseEmail(email: string): string {
    let [local, domain] = email.split('@')
    local = local.split('+')[0].replace(/\./g, '')
    return `${local}@${domain}`
  }

  return new Set(emails.map(parseEmail)).size
}

console.log(
  numUniqueEmails([
    'test.email+alex@leetcode.com',
    'test.e.mail+bob.cathy@leetcode.com',
    'testemail+david@lee.tcode.com',
  ])
)
// 如果在本地名称中添加加号（'+'），
// 则会忽略第一个加号后面的所有内容。
// 这允许过滤某些电子邮件，例如 m.y+name@email.com 将转发到 my@email.com。
// （同样，此规则不适用于域名。）
// 如果在电子邮件地址的本地名称部分中的某些字符之间添加句点（'.'），
// 则发往那里的邮件将会转发到本地名称中没有点的同一地址。
// 例如，"alice.z@leetcode.com” 和 “alicez@leetcode.com”
// 会转发到同一电子邮件地址。
// 我们会向列表中的每个地址发送一封电子邮件。实际收到邮件的不同地址有多少？
