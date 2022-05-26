  if (n === 0) return '0'

  const sb: number[] = []
  while (n > 0) {
    const [div, mod] = [Math.floor(n / 2), n % 2]
    sb.push(Math.abs(mod))
    n = -div
  }

  return sb.reverse().join('')