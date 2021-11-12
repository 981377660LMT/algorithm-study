function getExtname(fileName: string) {
  const dotIndex = fileName.lastIndexOf('.')
  if ([-1, 0].includes(dotIndex)) {
    return ''
  } else {
    return fileName.slice(dotIndex + 1)
  }
}

console.log(getExtname('asd'))
console.log(getExtname('.asd'))
console.log(getExtname('test.py'))
