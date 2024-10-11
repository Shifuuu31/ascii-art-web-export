package source

import (
	"strings"
)

// ascii art generator function
func GenerateAsciiArt(text string, fontMap map[rune][]string) string {
    var result []string
    lines := strings.Split(text, "\n")

    for _, line := range lines {
        var asciiLines [8]string
        for _, char := range line {
            if charArt, found := fontMap[char]; found {
                for i, asciiLine := range charArt {
                    asciiLines[i] += asciiLine
                }
            }
        }
        result = append(result, strings.Join(asciiLines[:], "\n"))
    }

    return strings.Join(result, "\n\n")
}
