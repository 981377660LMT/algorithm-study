  const count = Math.ceil(b.length / a.length)
  const str = a.repeat(count)
  return str.includes(b) ? count : (str + a).includes(b) ? count + 1 : -1