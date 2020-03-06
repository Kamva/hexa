package kitty

import "github.com/labstack/gommon/color"

// PrintBanner print the app's banner
func PrintBanner(banner, product, version, website string) {
	if banner == "" {
		printKittyBanner(product)
		return
	}
	c := color.New()
	c.Printf(banner, c.Yellow(product), c.Red("v"+version), c.Blue(website))
}

func printKittyBanner(product string) {
	PrintBanner(banner, product, Version, "")
}
