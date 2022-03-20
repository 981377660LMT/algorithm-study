// node 端： toBase64:
Buffer.from('123').toString('base64')
// 浏览器端： toBase64
window.btoa('123')
// node 端：decode:
Buffer.from('MTIz', 'base64').toString()
// 浏览器端：decode
window.atob('MTIz')
