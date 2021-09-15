const sss = new Set([1, 2, 3])
sss.forEach(item => console.log(item))
for (const item of sss.entries()) {
  console.log(item)
}
console.log(new Map([...sss].entries()))
