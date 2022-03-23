function groupBy(str: string): string {
  const groups: string[][] = []
  let pre = ''

  for (const char of str) {
    if (char !== pre) {
      const newGroup = [char]
      groups.push(newGroup)
    } else {
      groups[groups.length - 1].push(char)
    }

    pre = char
  }

  return groups.map(group => `${group[0]}${group.length}`).join('')
}

console.log(groupBy('abbcccdddd'))

export {}
