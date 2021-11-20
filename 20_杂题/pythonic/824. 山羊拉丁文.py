class Solution:
    def toGoatLatin(self, sentence: str) -> str:
        VOWEL = set('aeiouAEIOU')

        def transform(word, i):
            if word[0] not in VOWEL:  # 单词以辅音(consonant)字母开始
                word = word[1:] + word[0]
            return word + 'ma' + 'a' * i

        return ' '.join([transform(word, i) for i, word in enumerate(sentence.split(), 1)])

