// @ts-ignore
import MagicString from './MagicString.ts'

export default MagicString
// @ts-ignore
export { default as Bundle } from './Bundle.ts'
// @ts-ignore
export { default as SourceMap } from './SourceMap.ts'

if (require.main === module) {
  const s = new MagicString('hello world')
  console.log(s.toString())
}
