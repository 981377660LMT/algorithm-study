// node myScript.js -s --test --cool=true
hasFlags('-s') // true
hasFlags('--test', 'cool=true', '-s') // true
hasFlags('special') // false

function hasFlags(...flags: string[]) {
  return flags.every(flag => process.argv.includes(/^-{1,2}/.test(flag) ? flag : '--' + flag))
}
