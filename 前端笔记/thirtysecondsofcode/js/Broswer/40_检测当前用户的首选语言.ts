// Detects the preferred language of the current user.
// 检测当前用户的首选语言。:使用 NavigationLanguage.language 或第一个 NavigationLanguage.languages (如果可用的话) ，否则返回 defaultLang。
detectLanguage() // 'nl-NL'

function detectLanguage(defaultLang = 'en-US') {
  return (
    navigator.language ||
    (Array.isArray(navigator.languages) && navigator.languages[0]) ||
    defaultLang
  )
}
