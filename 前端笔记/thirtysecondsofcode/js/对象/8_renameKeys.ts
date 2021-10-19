export {}
// 用提供的值替换多个对象键的名称。
const obj = { name: 'Bobo', job: 'Front-End Master', shoeSize: 100 }
renameKeys({ name: 'firstName', job: 'passion' }, obj)
// { firstName: 'Bobo', passion: 'Front-End Master', shoeSize: 100 }

function renameKeys(keysMap: Record<string, string>, obj: Record<PropertyKey, any>) {
  return Object.keys(obj).reduce(
    (pre, key) => ({
      ...pre,
      ...{ [keysMap[key] ?? key]: obj[key] },
    }),
    {}
  )
}
