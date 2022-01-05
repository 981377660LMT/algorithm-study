const target = '*'
const sentence1 = 'asv*as'
const sentence2 = 'asvas'
const pos1 = sentence1.indexOf(target)
const pos2 = sentence2.indexOf(target)
console.log(sentence1.slice(0, pos1 + 1))
console.log(sentence2.slice(0, pos2 + 1))

export {}
