let line = readline()

const n = readline()

for (let i = 0; i < n; i++) {
  const line2 = readline().split(' ')
  const start = parseInt(line2[0])
  const len = parseInt(line2[1])
  const temp = line.slice(start, len).split('').reverse().join('')
  line = line.slice(0, start + len) + temp + line.slice(start + len)
}

print(line)
