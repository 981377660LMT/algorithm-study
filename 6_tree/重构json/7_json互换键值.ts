const userToSkill = {
  robert: ['programming', 'design', 'reactjs'],
  kimia: ['java', 'backend', 'services'],
  patrick: ['reactjs'],
  chris: ['reactjs', 'programming']
} as const

// 转换成
const skillToUser = {
  programming: ['robert', 'chris'],
  reactjs: ['patrick', 'robert', 'chris'],
  java: ['kimia'],
  backend: ['kimia'],
  services: ['kimia'],
  design: ['robert']
}

type ExtractTuplePropName<U extends Record<string, readonly string[]>> = {
  [K in keyof U]: U[K][number]
}[keyof U]

type SkillName = ExtractTuplePropName<typeof userToSkill>
type UserName = keyof typeof userToSkill

const res = new Map<SkillName, UserName[]>()
const getObjKeys = <O>(obj: O) => Object.keys(obj) as Array<keyof O>

getObjKeys(userToSkill).forEach(user => {
  userToSkill[user].forEach(skill => {
    !res.has(skill) && res.set(skill, [])
    res.get(skill)!.push(user)
  })
})

console.log(Object.fromEntries(res.entries()))
export {}
