import ErrorStackParser from 'error-stack-parser'

const frame = ErrorStackParser.parse(new Error('as'))
console.log(frame[0].lineNumber, frame[0].columnNumber)
