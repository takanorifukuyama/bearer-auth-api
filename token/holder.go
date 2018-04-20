package token

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"regexp"
)

const AUTH_TOKENS = "AUTH_TOKENS"

type Holder struct {
	tokens map[string][]*regexp.Regexp
	keys   []string
}

func NewHolder() *Holder {
	rawTokens := os.Getenv(AUTH_TOKENS)
	if len(rawTokens) == 0 {
		rawTokens = "{}"
	}

	tokens := map[string][]*regexp.Regexp{}
	keys := []string{}

	var obj interface{}
	if err := json.Unmarshal([]byte(rawTokens), &obj); err == nil {
		switch jsonTokens := obj.(type) {
		case map[string]interface{}:
			for k, v := range jsonTokens {
				rv := reflect.ValueOf(v)
				if rv.Kind() == reflect.Slice {
					sl := make([]*regexp.Regexp, 0, 0)
					for i := 0; i < rv.Len(); i++ {
						switch tokenValue := rv.Index(i).Interface().(type) {
						case string:
							tokenRe, err := regexp.Compile(tokenValue)
							if err == nil && tokenRe != nil {
								sl = append(sl, tokenRe)
							}
						}
					}
					tokens[k] = sl
				}
			}
		}
		for k, _ := range tokens {
			keys = append(keys, k)
		}
	} else {
		log.Printf("AUTH_TOKENS parse failed: %v\n", err)
	}

	return &Holder{
		tokens: tokens,
		keys:   keys,
	}
}

func (holder *Holder) GetTokens() []string {
	return holder.keys
}

func (holder *Holder) HasToken(token string) bool {
	_, ok := holder.tokens[token]
	return ok
}

func (holder *Holder) GetAllowedPaths(token string) []*regexp.Regexp {
	return holder.tokens[token]
}
