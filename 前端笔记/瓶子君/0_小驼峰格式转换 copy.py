class Solution:
    def format(self, name: str):
        # write code here
        if "_" not in name:
            return name[0].lower() + name[1:]
        words = [word.capitalize() for word in name.lower().split("_") if word != ""]
        words[0] = words[0].lower()
        return "".join(words)


print("a2weA".capitalize())
