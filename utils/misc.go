package utils

// could make this a method of builder potentially to prevent the double copy when making builder into string then reversing
func ReverseString(s string) string {
  runes := []rune(s)
  for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
      runes[i], runes[j] = runes[j], runes[i]
  }
  return string(runes)
}
