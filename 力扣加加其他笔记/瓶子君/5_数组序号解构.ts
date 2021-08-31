const arr = [1, 2, 3] as const
const { 0: firstA, 1: secA, 2: thirdA } = arr

console.log(firstA)

export {}
