package assets

import (
	"embed"
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/errors"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/tinne26/etxt/font"
	"golang.org/x/image/font/sfnt"
	"gopkg.in/yaml.v2"
)

//go:embed locales/**/*
var localeFS embed.FS

type Locale struct {
	LocaleCode string
	font       *sfnt.Font
	flag       ebiten.Image
	keyPairs   map[string]string
}

var availableLocales []string
var loadedLocales = make(map[string]*Locale)
var currentLocale *Locale

func InitLocales() {
	availableLocales = readLocaleDir()
	currentLocale = loadLocale(config.DEFAULT_LOCALE)
}

func readLocaleDir() []string {
	localeDir, err := localeFS.ReadDir(config.LOCALE_DIR)
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

	return locales
}

func CurrentLocale() string {
	return currentLocale.LocaleCode
}
func Locales() []string {
	return availableLocales
}

func loadLocale(locale string) *Locale {
	localePath := path.Join(config.LOCALE_DIR, locale)
	logger.Info("Loading locale %s", localePath)
	if _, err := localeFS.ReadDir(localePath); err != nil {
		logger.Error("Locale %s not found", locale)
		return nil
	}

	// Load flag image
	flagImg, _, err := ebitenutil.NewImageFromFileSystem(localeFS, path.Join(localePath, "flag.png"))
	if err != nil {
		return nil
	}

	// Load font by finding a .ttf in the locale directory
	// it can be named anything, but it must be a .ttf
	var localeFont *sfnt.Font
	fontPath := ""
	fontDir, err := localeFS.ReadDir(localePath)
	if err != nil {
		return nil
	}
	for _, entry := range fontDir {
		if !entry.IsDir() && path.Ext(entry.Name()) == ".ttf" {
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
	data, err := localeFS.ReadFile(path.Join(localePath, "strings.yaml"))
	if err != nil {
		return nil
	}

	var keyPairs map[string]string
	if err := yaml.Unmarshal(data, &keyPairs); err != nil {
		return nil
	}

	l := &Locale{
		LocaleCode: locale,
		font:       localeFont,
		flag:       *flagImg,
		keyPairs:   keyPairs,
	}
	validateLocale(l)

	loadedLocales[locale] = l
	return l
}

// Checks that the locale file has all keys,
// just emits warnings if keys are missing to help with development
func validateLocale(locale *Locale) {
	allKeys := types.AllLocaleKeys
	for _, key := range allKeys {
		if _, ok := locale.keyPairs[key]; !ok {
			fmt.Printf("Missing key %s in locale %s\n", key, locale.LocaleCode)
		}
	}
}

func getLocale(locale string) *Locale {
	for _, l := range availableLocales {
		if l == locale {
			if _, ok := loadedLocales[locale]; !ok {
				loadedLocales[locale] = loadLocale(locale)
			}
			return loadedLocales[locale]
		}
	}
	fmt.Println("Locale not found")
	return nil
}

func SetLocale(locale string) error {
	if currentLocale != nil && currentLocale.LocaleCode == locale {
		logger.Info("Locale already set to %s", locale)
		return nil
	}
	newLocale := getLocale(locale)
	if newLocale == nil {
		return errors.Raise(errors.UNKNOWN_LOCALE, locale)
	}

	currentLocale = newLocale
	return nil
}

func String(key string) string {
	if val, ok := currentLocale.keyPairs[key]; ok {
		return val
	}
	if config.FALLBACK_TO_DEFAULT_LOCALE {
		if val, ok := loadedLocales[config.DEFAULT_LOCALE].keyPairs[key]; ok {
			return val
		}
		return key
	}
	return key
}

func Flag() *ebiten.Image {
	return &currentLocale.flag
}

func Font() *sfnt.Font {
	return currentLocale.font
}
