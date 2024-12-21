package l

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/tinne26/etxt/font"
	"golang.org/x/image/font/sfnt"
	"gopkg.in/yaml.v2"
)

type Locale struct {
	LocaleCode string
	font       *sfnt.Font
	flag       ebiten.Image
	keyPairs   map[string]string
}

var currentLocale *Locale
var availableLocales = readLocaleDir()
var loadedLocales = make(map[string]*Locale)

func readLocaleDir() []string {
	localeDir, err := os.ReadDir(config.LOCALE_DIR)
	if err != nil {
		return nil
	}

	var locales []string
	for _, entry := range localeDir {
		if entry.IsDir() {
			locales = append(locales, entry.Name())
		}
	}

	return locales
}

func loadLocale(locale string) *Locale {
	localePath := filepath.Join(config.LOCALE_DIR, locale)
	if _, err := os.Stat(localePath); os.IsNotExist(err) {
		return nil
	}

	// Load flag image
	flagImg, _, err := ebitenutil.NewImageFromFile(filepath.Join(localePath, "flag.png"))
	if err != nil {
		return nil
	}
	// Resize flag image to 16 x 8
	// flagImg = flagImg.SubImage(image.Rect(0, 0, 16, 8)).(*ebiten.Image)

	// Load font by finding a .ttf in the locale directory
	// it can be named anything, but it must be a .ttf
	var localeFont *sfnt.Font
	fontPath := ""
	fontDir, err := os.ReadDir(localePath)
	if err != nil {
		return nil
	}
	for _, entry := range fontDir {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".ttf" {
			fontPath = filepath.Join(localePath, entry.Name())
			break
		}
	}
	if fontPath != "" {
		bytes, err := os.ReadFile(fontPath)
		if err == nil {
			f, _, err := font.ParseFromBytes(bytes)
			if err == nil {
				localeFont = f
			}
		}

	}

	// Load key pairs from JSON
	data, err := os.ReadFile(filepath.Join(localePath, "strings.yaml"))
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
// just emits warnings if keys are missing
func validateLocale(_ *Locale) {
	// TODO: Implement
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

func Change(locale string) {
	if currentLocale != nil && currentLocale.LocaleCode == locale {
		return
	}
	newLocale := getLocale(locale)
	if newLocale == nil {
		return
	}

	currentLocale = newLocale
	ui.SetFont(currentLocale.font)
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
