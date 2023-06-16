function hello(name) {
  // resolve 和 reject 向channel发送消息
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      if (name === 'fail') {
        reject(new Error('failed'))
      } else {
        resolve(`hello ${name}`)
      }
    }, 1e3)
  })
}

async function main() {
  try {
    let output = await hello('bob')
    console.log(output)

    output = await hello('fail')
    console.log(output)
  } catch (err) {
    console.log(err.message)
  }
}

main()
