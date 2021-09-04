// 几个api
declare module window {
  function btoa(str: string): string
  function atob(str: string): string
}

// node
console.log(Buffer.from('123').toString('base64'))
console.log(Buffer.from('MTIz', 'base64').toString())

//  浏览器
// console.log(window.btoa('123'))
// console.log(window.atob('MTIz'))
