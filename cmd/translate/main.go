package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	underarmour "github.com/blaskovicz/go-underarmour"
	gpx "github.com/ptrv/go-gpx"
)

func main() {
	token := flag.String("cookie-auth-token", "", "Cookie: auth-token=<value>")
	route := flag.Int("route", 0, "/routes/view/<value>")
	file := flag.String("file", "", "outfile.gpx")
	flag.Parse()
	if token != nil && *token != "" {
		os.Setenv("UNDERARMOUR_COOKIE_AUTH_TOKEN", *token)
	}
	if route == nil || *route <= 0 {
		panic("invalid route provided")
	}
	if file == nil || *file == "" {
		temp := fmt.Sprintf("./%d.gpx", *route)
		file = &temp
	}
	client := underarmour.New()
	rgpx, err := client.ReadRouteGPX(*route)
	if err != nil {
		panic(err)
	}
	garminGpx := rgpx.Clone()

	spoofTime := time.Now().UTC().Format(time.RFC3339Nano)
	garminGpx.Metadata = &gpx.Metadata{
		Links:     []gpx.Link{gpx.Link{Text: "Garmin Connect", URL: "connect.garmin.com"}},
		Timestamp: spoofTime,
	}
	garminGpx.Tracks[0].Type = "running"
	waypoints := garminGpx.Tracks[0].Segments[0].Waypoints
	for i, _ := range waypoints {
		waypoints[i].Timestamp = spoofTime
	}

	f, err := os.Create(*file)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(garminGpx.ToXML())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s written with garmin-adjusted route %d, now import to garmin and save as route!\n", f.Name(), *route)
}
