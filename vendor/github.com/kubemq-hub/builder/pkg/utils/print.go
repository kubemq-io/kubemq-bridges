package utils

import (
	"fmt"
	"github.com/gookit/color"
	"sort"
	"strings"
)

var themes = map[string]*color.Theme{
	"info":  color.Info,
	"note":  color.Note,
	"light": color.Light,
	"error": color.Error,

	"debug":   color.Debug,
	"danger":  color.Danger,
	"notice":  color.Notice,
	"success": color.Success,
	"comment": color.Comment,
	"primary": color.Primary,
	"warning": color.Warn,

	"question":  color.Question,
	"secondary": color.Secondary,
}

func Info(format string, args ...interface{}) {
	color.Info.Prompt(format, args...)
}
func Warn(format string, args ...interface{}) {
	color.Warn.Prompt(format, args...)
}
func Error(format string, args ...interface{}) {
	color.Error.Prompt(format, args...)
}

func Println(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	color.Print(fmt.Sprintf("%s\n", str))
}
func Print(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	color.Print(str)
}
func Block(theme string, format string, args ...interface{}) {
	colorTheme := themes[theme]
	colorTheme.Block(format, args...)
}

func stringSplitter(str string) string {
	lines := strings.Split(str, "\n")
	if len(lines) <= 1 {
		return str
	}
	newLines := []string{
		" |-",
	}
	for _, line := range lines {
		newLines = append(newLines, fmt.Sprintf("    %s", line))
	}
	return strings.Join(newLines, "\n")
}
func MapToYaml(val map[string]string) string {
	var list []string
	for key, val := range val {
		if val != "" {
			list = append(list, fmt.Sprintf("  <red>%s:</> %s", key, stringSplitter(val)))
		}
	}
	sort.Strings(list)
	return strings.Join(list, "\n")
}
func MapArrayToYaml(ml []map[string]string) string {
	var output []string

	for _, m := range ml {
		var list []string
		for key, val := range m {
			if val != "" {
				list = append(list, fmt.Sprintf("<red>%s:</> %s", key, stringSplitter(val)))
			}
		}
		sort.Strings(list)
		for i, line := range list {
			if i == 0 {
				output = append(output, fmt.Sprintf("- %s", line))
			} else {
				output = append(output, fmt.Sprintf("  %s", line))
			}
		}

	}
	return strings.Join(output, "\n")
}
