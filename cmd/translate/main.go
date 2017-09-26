package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	underarmour "github.com/blaskovicz/go-underarmour"
	"github.com/blaskovicz/mapmyrun-to-garmin/garmin"
)

func main() {
	garminPinM := flag.String("garmin-pin-m", "", "Cookie: pin.m=<value>")
	garminSession := flag.String("garmin-session", "", "Cookie: SESSION=<value>")
	uaToken := flag.String("underarmour-cookie-auth-token", "", "Cookie: auth-token=<value>")
	route := flag.Int("route", 0, "/routes/view/<value>")
	file := flag.String("file", "", "outfile.gpx")
	flag.Parse()
	if uaToken != nil && *uaToken != "" {
		os.Setenv("UNDERARMOUR_COOKIE_AUTH_TOKEN", *uaToken)
	}
	if garminPinM != nil && *garminPinM != "" {
		os.Setenv("GARMIN_COOKIE_PIN_M", *garminPinM)
	}
	if garminSession != nil && *garminSession != "" {
		os.Setenv("GARMIN_COOKIE_SESSION", *garminSession)
	}
	if route == nil || *route <= 0 {
		panic("invalid route provided")
	}
	if file == nil || *file == "" {
		temp := fmt.Sprintf("./%d.gpx", *route)
		file = &temp
	}
	gclient := garmin.New()

	// auth test
	cModel, err := gclient.ReadCourse(1)
	if err != nil {
		panic(fmt.Errorf("failed to read garmin course: %s", err))
	}

	uaclient := underarmour.New()
	uaRoute, err := uaclient.ReadRoute(*route)
	if err != nil {
		panic(fmt.Errorf("Failed to read route from mapmyrun: %s", err))
	}
	rgpx, err := uaclient.ReadRouteGPX(*route)
	if err != nil {
		panic(fmt.Errorf("Failed to read gpx route from mapmyrun: %s", err))
	}

	c, err := gclient.MassageCourseFromGPX(uaRoute.Distance, rgpx)
	if err != nil {
		panic(fmt.Errorf("Failed to massage garmin course from underarmour gpx: %s", err))
	}
	c.ActivityTypePk = cModel.ActivityTypePk
	public := 1
	running := 3
	// pid := 16226540
	c.RulePK = &public
	c.SourceTypeID = &running
	c.Name = uaRoute.Name
	c.Description = fmt.Sprintf(`MapMyRun route http://mapmyrun.com/routes/view/%d imported via https://github.com/blaskovicz/mapmyrun-to-garmin on %s.`, *route, time.Now().Format(time.RFC822))
	// cannot be set on create c.VirtualPartnerID = &pid

	// TODO allow using metadata to update course
	course, err := gclient.CreateCourse(c)
	if err != nil {
		panic(fmt.Errorf("Failed to create garmin course: %s", err))
	}
	fmt.Printf("Successfully imported course https://connect.garmin.com/modern/courses/%d '%s' into Garmin Connect!\n", *course.ID, course.Name)

	/*
		To just write it to a garmin activity style gpx
		file that can be imported / cloned to a route...

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
		fmt.Printf("%s written with garmin-adjusted route %d, now import to garmin and save as route!\n", f.Name(), *route)*/
}
