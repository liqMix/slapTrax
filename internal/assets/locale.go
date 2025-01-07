package assets

import (
	"embed"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/tinne26/etxt/font"
	"golang.org/x/image/font/sfnt"
	"gopkg.in/yaml.v2"
)

//go:embed locales/**/*
var localeFS embed.FS

const (
	DefaultLocaleCode = "en-us"

	localesDir     = "locales"
	fontExt        = ".ttf"
	localeFilename = "strings.yaml"
)

type Locale struct {
	LocaleCode string
	font       *sfnt.Font
	flag       ebiten.Image
	keyPairs   map[string]string
}

func (l *Locale) GetString(key string) (string, bool) {
	val, ok := l.keyPairs[key]
	if ok {
		return val, true
	}

	return "", false
}

var loadedLocales = make(map[string]*Locale)
var defaultLocale *Locale
var currentLocale *Locale
var availableLocales []string

func InitLocales(startingLocale string) {
	readLocaleDir()
	defaultL, err := loadLocale(DefaultLocaleCode)
	if err != nil {
		logger.Fatal("Failed to load default locale %s", DefaultLocaleCode)
	}

	defaultLocale = defaultL
	currentLocale = defaultLocale

	startL, err := loadLocale(startingLocale)
	if err != nil {
		logger.Error("Failed to load starting locale %s", startingLocale)
	} else {
		currentLocale = startL
	}
}

func readLocaleDir() []string {
	localeDir, err := localeFS.ReadDir(localesDir)
	if err != nil {
		return nil
	}

	var locales []string
	for _, entry := range localeDir {
		if entry.IsDir() {
			logger.Debug("Found locale %s", entry.Name())
			locales = append(locales, entry.Name())
		}
	}

	availableLocales = locales
	return locales
}

func CurrentLocale() string {
	return currentLocale.LocaleCode
}

func Locales() []string {
	return availableLocales
}

func loadLocale(locale string) (*Locale, error) {
	localePath := path.Join(localesDir, locale)
	logger.Info("Loading locale %s", localePath)
	if _, err := localeFS.ReadDir(localePath); err != nil {
		return nil, err
	}

	// Load flag image
	flagImg, _, err := ebitenutil.NewImageFromFileSystem(localeFS, path.Join(localePath, "flag.png"))
	if err != nil {
		return nil, err
	}

	// Load font by finding a .ttf in the locale directory
	// it can be named anything, but it must be a .ttf
	var localeFont *sfnt.Font
	fontPath := ""
	fontDir, err := localeFS.ReadDir(localePath)
	if err != nil {
		return nil, err
	}
	for _, entry := range fontDir {
		if !entry.IsDir() && path.Ext(entry.Name()) == fontExt {
			fontPath = path.Join(localePath, entry.Name())
			break
		}
	}
	if fontPath != "" {
		bytes, err := localeFS.ReadFile(fontPath)
		if err == nil {
			f, _, err := font.ParseFromBytes(bytes)
			if err == nil {
				localeFont = f
			}
		}

	}

	// Load key pairs from JSON
	data, err := localeFS.ReadFile(path.Join(localePath, localeFilename))
	if err != nil {
		return nil, err
	}

	var keyPairs map[string]string
	if err := yaml.Unmarshal(data, &keyPairs); err != nil {
		return nil, err
	}

	l := &Locale{
		LocaleCode: locale,
		font:       localeFont,
		flag:       *flagImg,
		keyPairs:   keyPairs,
	}
	validateLocale(l)

	loadedLocales[locale] = l
	return l, nil
}

// Checks that the locale file has all keys,
// just emits warnings if keys are missing to help with development
func validateLocale(locale *Locale) {
	if locale == nil || locale.LocaleCode == DefaultLocaleCode {
		return
	}

	logger.Info("Validating locale %s", locale.LocaleCode)
	allKeys := defaultLocale.keyPairs
	for _, key := range allKeys {
		if _, ok := locale.keyPairs[key]; !ok {
			logger.Warn("\tMissing key %s in locale %s\n", key, locale.LocaleCode)
		}
	}
}

func SetLocale(l string) error {
	if currentLocale.LocaleCode == l {
		logger.Warn("Locale already set to %s", l)
		return nil
	}

	var err error
	locale, ok := loadedLocales[l]
	if !ok {
		locale, err = loadLocale(l)
		if err != nil {
			return err
		}
		loadedLocales[l] = locale
	}

	currentLocale = locale
	ebiten.SetWindowTitle(GetLocaleString("title"))
	return nil
}

func GetLocaleString(key string) string {
	var val string
	var ok bool
	val, ok = currentLocale.GetString(key)
	if !ok {
		val, ok = defaultLocale.GetString(key)
		if !ok {
			val = key
		}
	}
	return val
}

func GetDefaultLocaleString(key string) (string, *sfnt.Font) {
	if val, ok := defaultLocale.GetString(key); ok {
		return val, defaultLocale.font
	}
	return "", defaultLocale.font
}

func Flag() *ebiten.Image {
	return &currentLocale.flag
}

func Font() *sfnt.Font {
	return currentLocale.font
}
