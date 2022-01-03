function bitLength(num: number) {
  return 32 - Math.clz32(num)
}

if (require.main === module) {
  console.log(bitLength(3))
  console.log(bitLength(4))
}

export { bitLength }
