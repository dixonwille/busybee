package util

import "strings"

//CleanHost will prepend "https://" if no protocole is found before the string.
func CleanHost(host string) string {
	if host == "" {
		return ""
	}
	if strings.Index(host, "http://") == 0 || strings.Index(host, "https://") == 0 {
		return host
	}
	return "https://" + host
}

//CleanMention will prepend "@" if it is not found before the string.
func CleanMention(mention string) string {
	if mention == "" {
		return ""
	}
	mentionRunes := []rune(mention)
	if mentionRunes[0] != '@' {
		mentionRunes = append([]rune{'@'}, mentionRunes...)
	}
	return string(mentionRunes)
}

//CleanQuest is used to cleanup questions to dispaly to user.
func CleanQuest(question string) string {
	if question == "" {
		return ""
	}
	question = strings.Trim(question, " ")
	questRunes := []rune(question)
	if questRunes[len(questRunes)-1] != '?' {
		questRunes = append(questRunes, '?')
	}
	questRunes = append(questRunes, ' ')
	return string(questRunes)
}
