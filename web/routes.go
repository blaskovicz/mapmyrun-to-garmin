package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	underarmour "github.com/blaskovicz/go-underarmour"
	"github.com/blaskovicz/mapmyrun-to-garmin/garmin"
	"github.com/sirupsen/logrus"
)

func IndexRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Misconfiguration"))
	}
}

type garminImportModel struct {
	GarminSession     string `json:"garminSession"`
	MapMyRunAuthToken string `json:"mapmyrunAuthToken"`
	MapMyRunRouteID   int    `json:"mapmyrunRouteId"`
}

func writeError(r *http.Request, w http.ResponseWriter, err error) error {
	return writeJSON(r, w, map[string]interface{}{"message": fmt.Sprintf("%s", err)}, http.StatusBadRequest)
}
func writeJSON(r *http.Request, w http.ResponseWriter, model interface{}, status int) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&model)
	rLogger := logrus.WithFields(logrus.Fields{"path": r.URL.Path, "status": status, "model": model})
	if err != nil {
		rLogger.WithError(err).Warn()
	} else {
		rLogger.Info()
	}
	return err
}

func ApiPostGarminImport(w http.ResponseWriter, r *http.Request) {
	var garminImport garminImportModel
	var err error
	if err = json.NewDecoder(r.Body).Decode(&garminImport); err != nil {
		writeError(r, w, err)
		return
	} else if garminImport.MapMyRunAuthToken == "" || garminImport.MapMyRunRouteID < 1 || garminImport.GarminSession == "" {
		writeError(r, w, fmt.Errorf("garminSession, mapmyrunAuthToken, and mapmyrunRouteId must all be set"))
		return
	}

	logrus.WithFields(logrus.Fields{"path": r.URL.Path, "map-my-run.route-id": garminImport.MapMyRunRouteID, "at": "start"}).Info()

	garminClient := garmin.New().SetCookieSession(garminImport.GarminSession)
	mapmyrunClient := underarmour.New().SetCookieAuthToken(garminImport.MapMyRunAuthToken)
	route, err := mapmyrunClient.ReadRoute(garminImport.MapMyRunRouteID)
	if err != nil {
		writeError(r, w, fmt.Errorf("Failed to read route %d info from mapmyrun: %s", garminImport.MapMyRunRouteID, err))
		return
	}
	gpxRoute, err := mapmyrunClient.ReadRouteGPX(garminImport.MapMyRunRouteID)
	if err != nil {
		writeError(r, w, fmt.Errorf("Failed to read route %d gpx data from mapmyrun: %s", garminImport.MapMyRunRouteID, err))
		return
	}
	garminCourse, err := garminClient.MassageCourseFromGPX(route.Distance, gpxRoute)
	if err != nil {
		writeError(r, w, fmt.Errorf("Failed to derrive garmin course from mapmyrun route %d data and gpx: %s", garminImport.MapMyRunRouteID, err))
		return
	}

	// hard-coded defaults that will need to change ... TODO
	public := 1
	running := 3
	garminCourse.ActivityTypePk = 1 // not sure what this is yet, but it's 1 for my routes
	garminCourse.RulePK = &public
	garminCourse.SourceTypeID = &running
	garminCourse.Name = route.Name
	garminCourse.Description =
		fmt.Sprintf(`MapMyRun route http://mapmyrun.com/routes/view/%d imported via https://mapmyrun-to-garmin.carlyzach.com (https://github.com/blaskovicz/mapmyrun-to-garmin) at %s.`, garminImport.MapMyRunRouteID, time.Now().Format(time.RFC822))

	// TODO allow using metadata to update course
	fullGarminCourse, err := garminClient.CreateCourse(garminCourse)
	if err != nil {
		writeError(r, w, fmt.Errorf("Failed to create course on garmin from mapmyrun route %d: %s", garminImport.MapMyRunRouteID, err))
		return
	}

	writeJSON(r, w, fullGarminCourse, http.StatusCreated)
}
