package hexa

import "github.com/labstack/gommon/color"

// PrintBanner print the app's banner
func PrintBanner(banner, product, version, website string) {
	if banner == "" {
		printDefaultBanner(product)
		return
	}
	c := color.New()
	c.Printf(banner, c.Yellow(product), c.Red("v"+version), c.Blue(website))
}

func printDefaultBanner(product string) {
	PrintBanner(banner, product, Version, "")
}
