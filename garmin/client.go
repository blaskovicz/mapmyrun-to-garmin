package garmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	geo "github.com/kellydunn/golang-geo"
	gpx "github.com/ptrv/go-gpx"
)

const (
	DefaultRootURI = "https://connect.garmin.com/modern/proxy"
	envVarPrefix   = "GARMIN"
)

type Client struct {
	rootURI string
	cookie  struct {
		session string
		pinM    string
	}
	accessToken string // TODO
}

func New() *Client {
	rootURI := os.Getenv(envVarPrefix + "_ROOT_URI")
	if rootURI == "" {
		rootURI = DefaultRootURI
	}
	cookieSession := os.Getenv(envVarPrefix + "_COOKIE_SESSION")
	cookiePinM := os.Getenv(envVarPrefix + "_COOKIE_PIN_M")
	c := &Client{
		rootURI: rootURI,
	}
	c.cookie.session = cookieSession
	c.cookie.pinM = cookiePinM
	return c
}

func (c *Client) uri(path string, pathArgs ...interface{}) string {
	return fmt.Sprintf("%s%s", c.rootURI, fmt.Sprintf(path, pathArgs...))
}

// do a request, return the undread response if no errors and 200 OK
func (c *Client) doWithResponse(req *http.Request) (*http.Response, error) {
	if c.cookie.session == "" {
		return nil, fmt.Errorf("missing cookie.session for request")
	}
	req.AddCookie(&http.Cookie{Name: "SESSION", Value: c.cookie.session})
	if c.cookie.pinM == "" {
		return nil, fmt.Errorf("missing cookie.pinM for request")
	}
	req.AddCookie(&http.Cookie{Name: "pin.m", Value: c.cookie.pinM})
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("NK", "NT")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	// error
	if resp.StatusCode != http.StatusOK {
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s response payload: %s", resp.Status, err)
		}
		return nil, fmt.Errorf("request failed with %s: %s", resp.Status, string(rawBody))
	}
	return resp, nil
}

func (c *Client) do(req *http.Request, decodeTarget interface{}) error {
	if decodeTarget != nil {
		if decodeKind := reflect.TypeOf(decodeTarget).Kind(); decodeKind != reflect.Ptr {
			return fmt.Errorf("invalid decode target type %s (need %s)", decodeKind.String(), reflect.Ptr.String())
		}
	}
	resp, err := c.doWithResponse(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if decodeTarget != nil {
		err = json.NewDecoder(resp.Body).Decode(decodeTarget)
		if err != nil {
			return fmt.Errorf("failed to decode payload: %s", err)
		}
	}
	return nil
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Point struct {
	Timestamp *int64   `json:"timestamp,omitempty"`
	Distance  *float64 `json:"distance"`
	Elevation float64  `json:"elevation"`
	Coordinate
}
type Line struct {
	DistanceInMeters float64 `json:"distanceInMeters"`
	SortOrder        int     `json:"sortOrder"`
	//Points           []Point `json:"points"`
	NumberOfPoints int     `json:"numberOfPoints"`
	Bearing        float32 `json:"bearing"`
	GeoRoutePk     *int    `json:"geoRoutePk"`
}
type Course struct {
	ActivityTypePk       int      `json:"activityTypePk"`
	ID                   *int     `json:"courseId"`
	VirtualPartnerID     *int     `json:"virtualPartnerId"`
	UserProfilePk        *int     `json:"userProfilePk"`
	GeoRoutePk           *int     `json:"geoRoutePk"`
	RulePK               *int     `json:"rulePK"` // 1 for public, 2 for private
	SourceTypeID         *int     `json:"sourceTypeId"`
	Name                 string   `json:"courseName"`
	Description          string   `json:"description"`
	OpenStreetMap        bool     `json:"openStreetMap"`
	MatchedToSegments    bool     `json:"matchedToSegments"`
	FirstName            string   `json:"firstName"`
	Lastname             string   `json:"lastName"`
	DisplayName          string   `json:"displayName"`
	ElevationGain        float64  `json:"elevationGainMeter"`
	ElevationLoss        float64  `json:"elevationLossMeter"`
	Distance             float64  `json:"distanceMeter"`
	IncludeLaps          bool     `json:"includeLaps"`
	ElapsedSeconds       *float64 `json:"elapsedSeconds"`
	SpeedMetersPerSecond *float64 `json:"speedMeterPerSecond"`
	StartPoint           Point    `json:"startPoint"`
	GeoPoints            []Point  `json:"geoPoints"`
	CoursePoints         []Point  `json:"coursePoints"`
	BoundingBox          struct {
		Center              *Coordinate `json:"center"`
		LowerLeft           *Coordinate `json:"lowerLeft"`
		UpperRight          *Coordinate `json:"upperRight"`
		LowerLeftLatIsSet   bool        `json:"lowerLeftLatIsSet"`
		LowerLeftLongIsSet  bool        `json:"lowerLeftLongIsSet"`
		UpperRightLatIsSet  bool        `json:"upperRightLatIsSet"`
		UpperRightLongIsSet bool        `json:"upperRightLongIsSet"`
	} `json:"boundingBox"`
	CreateDate string `json:"createDate"`
	UpdateDate string `json:"updateDate"`
	//CreateDate  time.Time `json:"createDate"` // non-standard time 2017-08-12T02:47:01.0
	//UpdateDate  time.Time `json:"updateDate"` // non-standard time 2017-08-12T02:47:01.0
	CourseLines []Line `json:"courseLines"`
}

func (c *Client) ReadCourse(courseID int) (*Course, error) {
	req, err := http.NewRequest("GET", c.uri("/course-service/course/%d/", courseID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	var r Course
	if err = c.do(req, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// [ [lat, lon, ele] ... ] -> [ [ lat, lon, ele ] ... ]
func (c *Client) ReadElevation(data [][]*float64) ([][]*float64, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal elevationrequest data: %s", err)
	}
	req, err := http.NewRequest("POST", c.uri("/course-service/course/elevation"), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	var newData [][]*float64
	if err = c.do(req, &newData); err != nil {
		return nil, err
	}
	return newData, nil
}

func (c *Client) CreateCourse(course *Course) (*Course, error) {
	payload, err := json.Marshal(course)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal course: %s", err)
	}
	req, err := http.NewRequest("POST", c.uri("/course-service/course/"), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	var r Course
	if err = c.do(req, &r); err != nil {
		return nil, err
	}
	for i, _ := range course.CourseLines {
		course.CourseLines[i].GeoRoutePk = r.GeoRoutePk
	}
	l, err := c.createSegments(course.CourseLines)
	if err != nil {
		return nil, err
	}
	r.CourseLines = l
	return &r, nil
}

func (c *Client) createSegments(segments []Line) ([]Line, error) {
	payload, err := json.Marshal(segments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal course lines: %s", err)
	}
	req, err := http.NewRequest("POST", c.uri("/course-service/course/routeSegments"), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	var s []Line
	if err = c.do(req, &s); err != nil {
		return nil, err
	}
	return s, nil
}

func (c *Client) MassageCourseFromGPX(distanceInMeters float64, course *gpx.Gpx) (*Course, error) {
	target := &Course{GeoPoints: []Point{}, Distance: distanceInMeters, CourseLines: []Line{}}
	waypoints := course.Tracks[0].Segments[0].Waypoints
	var lastPoint *Point
	var distanceTraveled float64
	var sortOffset int = 1
	var ele [][]*float64 = [][]*float64{}
	var pointEvery int = 1
	var pointsAccrued int
	var wpc = len(waypoints)

	// TODO better algorithm (distance?) here.
	// we need this because garmin ui / backend will choke if we try to render too many user points
	if wpc > 20 {
		// break up point every so the line map
		// in garmin connect isn't so heavy
		pointEvery = wpc / 20
	}

	for i, p := range waypoints {
		var point Point
		point.Latitude = p.Lat
		point.Longitude = p.Lon
		point.Elevation = p.Ele
		ele = append(ele, []*float64{&point.Latitude, &point.Longitude, &point.Elevation})
		if lastPoint == nil {
			lastPoint = &point
		} else {
			p1 := geo.NewPoint(lastPoint.Latitude, lastPoint.Longitude)
			p2 := geo.NewPoint(point.Latitude, point.Longitude)
			legDistance := 1000.0 * p1.GreatCircleDistance(p2)
			distanceTraveled += legDistance
			var newDist = distanceTraveled
			point.Distance = &newDist
			pointsAccrued++
			if i%pointEvery == 0 || wpc-1 == i {
				target.CourseLines = append(target.CourseLines, Line{SortOrder: sortOffset, DistanceInMeters: newDist, NumberOfPoints: pointsAccrued})
				sortOffset++
				pointsAccrued = 0
			}
			lastPoint = &point
		}
		// TODO point.Timestamp = p.Time
		target.GeoPoints = append(target.GeoPoints, point)

		if target.BoundingBox.LowerLeft == nil || (target.BoundingBox.LowerLeft.Latitude >= point.Latitude && target.BoundingBox.LowerLeft.Longitude >= point.Longitude) {
			target.BoundingBox.LowerLeft = &point.Coordinate
			target.BoundingBox.LowerLeftLatIsSet = true
			target.BoundingBox.LowerLeftLongIsSet = true
		}
		if target.BoundingBox.UpperRight == nil || (target.BoundingBox.UpperRight.Latitude <= point.Latitude && target.BoundingBox.UpperRight.Longitude <= point.Longitude) {
			target.BoundingBox.UpperRight = &point.Coordinate
			target.BoundingBox.UpperRightLatIsSet = true
			target.BoundingBox.UpperRightLongIsSet = true
		}
	}
	eleData, err := c.ReadElevation(ele)
	if err != nil {
		return nil, err
	}
	for i, _ := range eleData {
		if eleData[i][2] == nil {
			continue
		}
		target.GeoPoints[i].Elevation = *eleData[i][2]
	}
	target.StartPoint = target.GeoPoints[0]
	// these cannot be the same time TODO target.CoursePoints = target.GeoPoints
	return target, nil
}
