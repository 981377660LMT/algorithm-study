console.log(0 ?? 'ok') // 只有当前面为undefined或null时才去后面的值
console.log(0 || 'ok') // 前面为falsy都取后面的值
console.log(undefined ?? 'ok')
console.log(undefined || 'ok')
