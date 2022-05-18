const str = '我123的123456789的123456987'

console.log(
  str.replace(/(我\d{3})(的)(\d{9})(的)(\d{9})/g, (...args) => {
    return `${args[1]}${args[2]}${args[3]}${args[4]}${args[5]}`
  })
)

// console.log(str.match(/的(\d)+/g))
