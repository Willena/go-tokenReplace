package main

import (
	"errors"
	"regexp"
	"strings"
)

var tokenFinderRegex = regexp.MustCompile(`(?m)(\'|\")*(\$([a-zA-Z0-9\._\-]*))(\'|\")*`)

type ReplacementToken interface {
	GetReplacementValue(fail ...bool) (string, error)
}

type RawTokenValue struct {
	raw string
}

func (t *RawTokenValue) GetReplacementValue(fail ...bool) (string, error) {
	return t.raw, nil
}

func CreateRawTokenValue(str string) *RawTokenValue {
	return &RawTokenValue{raw: str}
}

type SanitizedTokenValue struct {
	value string
}

func (t *SanitizedTokenValue) GetReplacementValue(fail ...bool) (string, error) {
	return "\"" + strings.ReplaceAll(t.value, "\"", "\\\"") + "\"", nil
}

func CreateSanitizedValue(str string) *SanitizedTokenValue {
	return &SanitizedTokenValue{value: str}
}

type CompoundTokenValue struct {
	tokens             map[string]ReplacementToken
	str                string
	sanitized          bool
	failOnMissingParam bool
}

func (cToken *CompoundTokenValue) GetReplacementValue(fail ...bool) (string, error) {

	shouldFail := cToken.failOnMissingParam
	if len(fail) != 0 {
		shouldFail = fail[0]
	}

	tempStr := cToken.str
	other := ""

	for {
		match := tokenFinderRegex.FindStringIndex(tempStr)
		if match == nil {
			break
		}

		matchGroupsSub := tokenFinderRegex.FindStringSubmatch(tempStr)

		if value, ok := cToken.tokens[matchGroupsSub[3]]; ok {
			strValue, err := value.GetReplacementValue(shouldFail)
			if err != nil {
				return "", err
			}
			other += tempStr[0:match[0]] + strValue
		} else {
			other += tempStr[0:match[1]]
			if shouldFail {
				return "", errors.New("The token $" + matchGroupsSub[3] + " could not be bound (not found in the param table)")
			}
		}

		tempStr = tempStr[match[1]:]

	}
	other += tempStr
	if cToken.sanitized {
		return CreateSanitizedValue(other).GetReplacementValue()
	} else {
		return other, nil
	}
}

func (cToken *CompoundTokenValue) getSanitizedReplacementValue() (string, error) {
	old := cToken.sanitized
	cToken.sanitized = true
	out, err := cToken.GetReplacementValue()
	cToken.sanitized = old
	return out, err
}

func (cToken *CompoundTokenValue) WithFailures() *CompoundTokenValue {
	cToken.failOnMissingParam = true
	return cToken
}

func (cToken *CompoundTokenValue) Put(key string, value ReplacementToken) *CompoundTokenValue {
	cToken.tokens[key] = value
	return cToken
}

func (cToken *CompoundTokenValue) PutString(key string, value string) *CompoundTokenValue {
	cToken.tokens[key] = CreateSanitizedValue(value)
	return cToken
}

func (cToken *CompoundTokenValue) PutRaw(key string, value string) *CompoundTokenValue {
	cToken.tokens[key] = CreateRawTokenValue(value)
	return cToken
}

func CreateCompound(str string) *CompoundTokenValue {
	return &CompoundTokenValue{
		tokens:    make(map[string]ReplacementToken),
		str:       str,
		sanitized: false,
	}
}

func CreateSanitizedCompound(str string) *CompoundTokenValue {
	return &CompoundTokenValue{
		tokens:    make(map[string]ReplacementToken),
		str:       str,
		sanitized: true,
	}
}