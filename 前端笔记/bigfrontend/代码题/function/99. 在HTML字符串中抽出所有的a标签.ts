/**
 * @param {string} str
 * @return {string[]}
 */
function extract(str: string): string[] {
  // your code here
  return str.match(/(<a(\s[^>]*)?>)(.*?)(<\s*?\/\s*?a>)/g) || []
}

console.log(
  extract(`
<div>
    <a>link1< / a><a href="https://bfe.dev">link1< / a>
    <div<abbr>bfe</abbr>div>
    <div>
<abbr>bfe</abbr><a href="https://bfe.dev" class="link2"> <abbr>bfe</abbr>   <span class="l">l</span><span  class="i">i</span>   nk2   </a>
    </div>
</div>
`)
)
