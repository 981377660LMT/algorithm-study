const len = input()

for (let _ = 0; _ < len; _++) {
  const line = input()
  console.log(getAns(line))
}

function getAns(s) {
  if (s.length < 2) return 'Wrong'
  const first = s.charCodeAt(0)
  if (first < 65 || first > 122 || (first > 90 && first < 97)) return 'Wrong'

  const reg = new RegExp(/^[0-9]*$/)
  let num = 1
  for (let i = 1; i < s.length; i++) {
    if (!reg.test(s[i])) {
      let x = s.charCodeAt(i)
      if (x < 65 || x > 122 || (x > 90 && x < 97)) return 'Wrong'
    } else {
      num = 2
    }
  }

  return num == 2 ? 'Accept' : 'Wrong'
}
