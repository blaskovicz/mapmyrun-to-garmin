package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	underarmour "github.com/blaskovicz/go-underarmour"
	"github.com/blaskovicz/mapmyrun-to-garmin/garmin"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/routes/new", http.StatusFound)
}

// quick-n-dirty
var t = template.Must(template.New("route_form.tmpl").Parse(`
<html>
<head>
	<title>MapMyRun to Garmin</title>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css" integrity="sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ" crossorigin="anonymous">
	<script src="https://code.jquery.com/jquery-3.1.1.slim.min.js" integrity="sha384-A7FZj7v+d/sdmMqp/nOQwliLvUsJfDHW+k9Omg/a/EheAdgtzNs3hpfag6Ed950n" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/tether/1.4.0/js/tether.min.js" integrity="sha384-DztdAPBWPRXSA/3eYEEUWrWCy7G5KFbe8fFjk5JAIxUYHKkDx6Qin1DkWx51bBrb" crossorigin="anonymous"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/js/bootstrap.min.js" integrity="sha384-vBWWzlZJ8ea9aCX4pEW3rVHjgjt7zpkNpZk+02D9phzyeVkE+jo0ieGizqPLForn" crossorigin="anonymous"></script>
	<style>label.bold{font-weight:bold}</style>
	<script type="text/javascript"> var infolinks_pid = 3108833; var infolinks_wsid = 0; </script>
	<script type="text/javascript" src="//resources.infolinks.com/js/infolinks_main.js"></script>
</head>
<body>
	<div class="jumbotron">
		<h1><a href="/" style="color:inherit">MapMyRun to Garmin</a></h1>
		<p class="lead">Import <a href="http://mapmyrun.com" target="_blank" rel="noopener noreferrer">MapMyRun</a> Routes to <a href="http://connect.garmin.com" target="_blank" rel="noopener noreferrer">Garmin Connect</a></p>
		<hr class="my-1">
	</div>
	<div class="container">
		{{ if .flashError }}
			<div class="row">
				<div class="alert alert-danger col-12">
					{{.flashError}}
				</div>
			</div>
		{{ end }}
		{{ if .flashSuccess }}
			<div class="row">
				<div class="alert alert-success col-12">
					{{.flashSuccess}}
				</div>
			</div>
		{{ end }}
		<div class="row">
			<div class="col-12">
				<form method="POST" accept-charset="UTF-8">
					<div class="form-group {{ if .garminSessionE }} has-danger {{ end }}">
						<label for="garminSession" class="bold">Garmin <code>SESSIONID</code> Cookie</label>
						<input class="form-control" type="text" name="garminSession" id="garminSession" value="{{.garminSession}}">
						{{ if .garminSessionE }}<div class="form-control-feedback">It looks like this isn't field is invalid.</div>{{ end }}
						<small class="form-text text-muted">Find this on <a rel="noopener noreferrer" target="_blank" href="http://connect.garmin.com">Garmin Connect</a> under domain <i>connect.garmin.com</i> (eg: <code>59123-...</code>).</small>
					</div>
					<div class="form-group {{ if .mapmyrunAuthTokenE }} has-danger {{ end }}">
						<label for="mapmyrunAuthToken" class="bold">MapMyRun <code>auth-token</code> Cookie</label>
						<input class="form-control" type="text" name="mapmyrunAuthToken" id="mapmyrunAuthToken" value="{{.mapmyrunAuthToken}}">
						{{ if .mapmyrunAuthTokenE }}<div class="form-control-feedback">It looks like this isn't field is invalid.</div>{{ end }}
						<small class="form-text text-muted">Find this on <a rel="noopener noreferrer" target="_blank" href="http://mapmyrun.com">MapMyRun</a> under domain <i>mapmyrun.com</i> (eg: <code>US.123...</code>).</small>
					</div>
					<div class="form-group {{ if .mapmyrunIDE }} has-danger {{ end }}">
						<label for="mapmyrunIDE" class="bold">MapMyRun Route ID</label>
						<input class="form-control" type="text" name="mapmyrunID" id="mapmyrunID" value="{{.mapmyrunID}}">
						{{ if .mapmyrunIDE }}<div class="form-control-feedback">It looks like this isn't field is invalid.</div>{{ end }}
						<small class="form-text text-muted">Find this on <a rel="noopener noreferrer" target="_blank" href="https://mapmyrun.com/routes/my_routes">MapMyRun</a> (eg: http://mapmyrun.com/routes/view/123 would have ID <code>123</code>).</small>
					</div>
					{{ .csrfField }}
					<input type="submit" value="Import" class="btn btn-primary">
				</form>
			</div>
		</div>
		<div class="row">
			<div class="col-4">
				<p class="text-muted">Questions or Concerns? Check out my <a rel="noopener noreferrer" href="https://github.com/blaskovicz/mapmyrun-to-garmin" target="_blank">Github</a> repo.
				</p>
			</div>
		</div>
	</div>
</body>
</html>
`))

func NewRouteForm(ss sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := ss.Get(r, "session")
		if isRedirect := session.Values["isRedirect"]; isRedirect == nil {
			delete(session.Values, "flashError")
			delete(session.Values, "flashSuccess")
			delete(session.Values, "mapmyrunIDE")
			delete(session.Values, "mapmyrunAuthTokenE")
			delete(session.Values, "garminSessionE")
		}
		delete(session.Values, "isRedirect")
		session.Save(r, w)
		t.ExecuteTemplate(w, "route_form.tmpl", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
			"flashError":     session.Values["flashError"],
			"flashSuccess":   session.Values["flashSuccess"],

			"mapmyrunAuthToken":  session.Values["mapmyrunAuthToken"],
			"mapmyrunAuthTokenE": session.Values["mapmyrunAuthTokenE"],
			"mapmyrunID":         session.Values["mapmyrunID"],
			"mapmyrunIDE":        session.Values["mapmyrunIDE"],
			"garminSession":      session.Values["garminSession"],
			"garminSessionE":     session.Values["garminSessionE"],
		})
	}
}

func PostRouteForm(ss sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := ss.Get(r, "session")
		var hadError bool
		defer func() {
			if flashError := session.Values["flashError"]; flashError != nil {
				msg, _ := flashError.(string)
				fmt.Printf("[post-route-form.error] %s\n", msg)
			} else if hadError {
				session.Values["flashError"] = "Hmm, that doesn't look quite right..."
			}
			if flashSuccess := session.Values["flashSuccess"]; flashSuccess != nil {
				msg, _ := flashSuccess.(string)
				fmt.Printf("[post-route-form.success] %s\n", msg)
			}
			session.Values["isRedirect"] = true
			session.Save(r, w)
			http.Redirect(w, r, "/routes/new", http.StatusSeeOther)
		}()

		delete(session.Values, "flashError")
		delete(session.Values, "flashSuccess")
		delete(session.Values, "mapmyrunIDE")
		delete(session.Values, "mapmyrunAuthTokenE")
		delete(session.Values, "garminSessionE")

		mapmyrunID := r.PostForm.Get("mapmyrunID")
		mapmyrunRouteID, err := strconv.Atoi(mapmyrunID)
		session.Values["mapmyrunID"] = mapmyrunID
		if err != nil {
			hadError = true
			session.Values["mapmyrunIDE"] = true
		}

		mapmyrunAuthToken := r.PostForm.Get("mapmyrunAuthToken")
		session.Values["mapmyrunAuthToken"] = mapmyrunAuthToken
		if mapmyrunAuthToken == "" {
			hadError = true
			session.Values["mapmyrunAuthTokenE"] = true
		}

		garminSession := r.PostForm.Get("garminSession")
		session.Values["garminSession"] = garminSession
		if garminSession == "" {
			hadError = true
			session.Values["garminSessionE"] = true
		}

		if hadError {
			return
		}

		fmt.Printf("[post-route-form.info] Attempting import of mapmyrun.route=%d\n", mapmyrunRouteID)

		garminClient := garmin.New().SetCookieSession(garminSession)
		mapmyrunClient := underarmour.New().SetCookieAuthToken(mapmyrunAuthToken)
		route, err := mapmyrunClient.ReadRoute(mapmyrunRouteID)
		if err != nil {
			session.Values["flashError"] = fmt.Sprintf("Failed to read route %d from mapmyrun: %s", mapmyrunRouteID, err)
			return
		}
		gpxRoute, err := mapmyrunClient.ReadRouteGPX(mapmyrunRouteID)
		if err != nil {
			session.Values["flashError"] = fmt.Sprintf("Failed to read gpx data from mapmyrun: %s", err)
			return
		}
		garminCourse, err := garminClient.MassageCourseFromGPX(route.Distance, gpxRoute)
		if err != nil {
			session.Values["flashError"] = fmt.Sprintf("Failed to derrive garmin course from mapmyrun gpx: %s", err)
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
			fmt.Sprintf(`MapMyRun route http://mapmyrun.com/routes/view/%d imported via https://mapmyrun-to-garmin.herokuapp.com (https://github.com/blaskovicz/mapmyrun-to-garmin) at %s.`, mapmyrunRouteID, time.Now().Format(time.RFC822))

		// TODO allow using metadata to update course
		fullGarminCourse, err := garminClient.CreateCourse(garminCourse)
		if err != nil {
			session.Values["flashError"] = fmt.Sprintf("Failed to create course on garmin: %s", err)
			return
		}

		delete(session.Values, "mapmyrunID")
		session.Values["flashSuccess"] = fmt.Sprintf("Successfully imported course https://connect.garmin.com/modern/course/%d '%s' into Garmin Connect from http://mapmyrun.com/routes/view/%d. Don't forget to send the course to your device via Garmin Connect!", *fullGarminCourse.ID, fullGarminCourse.Name, mapmyrunRouteID)
	}
}
