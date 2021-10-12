function smallestUniqueSubstr(str: string) {
  const stack: string[] = []
  const visited = new Set<string>()
  const lastIndex = new Map<string, number>()

  for (let i = 0; i < str.length; i++) {
    lastIndex.set(str[i], i)
  }

  for (let i = 0; i < str.length; i++) {
    const char = str[i]

    if (!visited.has(char)) {
      while (
        stack.length &&
        stack[stack.length - 1] > char &&
        lastIndex.get(stack[stack.length - 1])! > i // 尚有存余
      ) {
        visited.delete(stack.pop()!)
      }

      visited.add(char)
      stack.push(char)
    }
  }

  return stack.join('')
}

console.log(smallestUniqueSubstr('xyzabcxyzabc'))
