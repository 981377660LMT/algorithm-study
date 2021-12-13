function interpret(command: string): string {
  return command.replace(/\(\)/g, 'o').replace(/\(al\)/g, 'al')

  const sb: string[] = []
  let bufferCount = 0

  for (const char of command) {
    if (char === 'G') {
      sb.push('G')
    } else if (char === ')') {
      if (bufferCount === 1) sb.push('o')
      else sb.push('al')
      bufferCount = 0
    } else {
      bufferCount++
    }
  }

  return sb.join('')
}

console.log(interpret('G()(al)'))
console.log(interpret('G()()()()(al)'))
