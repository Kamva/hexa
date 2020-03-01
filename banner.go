package kitty

import "github.com/labstack/gommon/color"

// PrintBanner print the app's banner
func PrintBanner(banner, product, version, website string) {
	if banner == "" {
		printKittyBanner()
		return
	}
	c := color.New()
	c.Printf(banner, c.Red("v"+Version), c.Blue(website))
}

func printKittyBanner() {
	PrintBanner(banner, "", Version, "")
}
