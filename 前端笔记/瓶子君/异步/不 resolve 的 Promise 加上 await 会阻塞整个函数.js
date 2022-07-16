const demo = () => new Promise((resolve, reject) => console.log(1))

async function main() {
  await demo()
  console.log(11)
}

console.log(1)
