package main

import (
	"auto-NewLine/env"
	"errors"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	exep, err := env.GetExecDir()
	if err != nil {
		println(err)
		return
	}

	app := &cli.App{
		Name:            "auto-NewLine",
		Usage:           Version,
		Description:     "指定したテキストファイルを最大文字数に合わせていい感じに改行します",
		Version:         Version,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "text",
				Aliases:  []string{"t"},
				Usage:    "改行するテキストファイルのパス",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "maxlen",
				Aliases: []string{"l"},
				Value:   30,
				Usage:   "最大文字数",
			},
			&cli.StringFlag{
				Name:    "encode",
				Aliases: []string{"e"},
				Value:   "shift-jis",
				Usage:   "エンコード",
			},
			&cli.PathFlag{
				Name:    "setting",
				Aliases: []string{"s"},
				Value:   path.Join(exep, "setting.yml"),
				Usage:   "設定ファイルのパス",
			},
		},
		Action: appfunc,
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func appfunc(c *cli.Context) error {
	maxlen := c.Int("maxlen")
	textp := c.Path("text")

	enc, err := GetEncoding(c.String("encode"))
	if err != nil {
		return err
	}

	text, err := LoadText(textp, enc)
	if err != nil {
		return err
	}
	if len([]rune(string(text))) < maxlen {
		return nil
	}

	stg, err := LoadSetting(c.Path("setting"))
	if err != nil {
		return err
	}

	minstrs, err := GetLineBreakable(string(text), *stg)
	if err != nil {
		return err
	}

	return WriteTextFile(textp, WithLimitJoinStrings(len([]rune(string(text))), minstrs, stg.Newline, maxlen), enc)
}

func GetEncoding(enc string) (encoding.Encoding, error) {
	e, ok := map[string]encoding.Encoding{
		"utf-8":     nil,
		"shift-jis": japanese.ShiftJIS,
	}[enc]

	if !ok {
		return nil, errors.New("the encoding was not found")
	}
	return e, nil
}

func LoadText(textp string, enc encoding.Encoding) ([]byte, error) {
	text, err := os.ReadFile(textp)
	if err != nil {
		return nil, err
	}

	if enc != nil {
		text, _, err = transform.Bytes(enc.NewDecoder(), text)
		if err != nil {
			return nil, err
		}
	}

	return text, nil
}

func LoadSetting(path string) (*Setting, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := Setting{}
	err = yaml.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func WithLimitJoinStrings(textlen int, strs []string, sep string, maxlen int) string {
	lines := []string{}
	line := strs[0]
	minlen := textlen / int(math.Ceil(float64(textlen)/float64(maxlen)))

	for _, str := range strs[1:] {
		linelen := len([]rune(line))
		strlen := len([]rune(str))

		if minlen <= linelen || linelen+strlen >= maxlen {
			lines = append(lines, line)
			line = str
			continue
		}

		line += str
	}

	if line != "" {
		lines = append(lines, line)
	}

	return strings.Join(lines, sep)
}

func GetLineBreakable(text string, stg Setting) ([]string, error) {
	ter, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	tokens := ter.Tokenize(text)
	lbrk := []string{}

	backfetrs := []string{}
	back := ""
	for i, token := range tokens {
		next := token.Surface
		nextfetrs := token.Features()[0:6]

		if i > 0 {
			canlbrk, err := CanLineBreak(back, next, backfetrs, nextfetrs, stg)
			if err != nil {
				return nil, err
			}
			if canlbrk {
				lbrk[len(lbrk)-1] += next
				back = next
				backfetrs = nextfetrs
				continue
			}
		}

		back = next
		backfetrs = nextfetrs
		lbrk = append(lbrk, next)
	}

	return lbrk, nil
}

func CanLineBreak(back string, next string, backfetrs []string, nextfetrs []string, stg Setting) (bool, error) {
	for _, backfetr := range backfetrs {
		for _, nextfetr := range nextfetrs {
			for _, ptnsdict := range stg.NotNLFeatures {
				ismatch, err := MatchKeyValue(backfetr, nextfetr, ptnsdict)
				if err != nil || ismatch {
					return ismatch, err
				}
			}
		}
	}

	for _, ptnsdict := range stg.NotNLStrings {
		ismatch, err := MatchKeyValue(back, next, ptnsdict)
		if err != nil || ismatch {
			return ismatch, err
		}
	}

	return false, nil
}

func MatchKeyValue(key string, value string, dict map[string]string) (bool, error) {
	for k, v := range dict {
		isbackm, err := regexp.MatchString(k, key)
		if err != nil {
			return false, err
		}

		isnextm, err := regexp.MatchString(v, value)
		if err != nil {
			return false, err
		}

		if isbackm && isnextm {
			return true, nil
		}
	}

	return false, nil
}

func WriteTextFile(path string, text string, enc encoding.Encoding) error {
	bytes := []byte(text)

	if enc != nil {
		b, _, err := transform.Bytes(enc.NewEncoder(), []byte(text))
		bytes = b
		if err != nil {
			return err
		}
	}

	return os.WriteFile(path, bytes, 0777)
}
