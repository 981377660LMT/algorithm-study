// atob() 可以用来解码通过Base64编码的字符串。
function myAtob(encoded: string): string {}

console.log(myAtob('QkZFLmRldg==')) // 'BFE.dev'
console.log(myAtob('Q')) // Error
