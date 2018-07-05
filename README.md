# mapmyrun-to-garmin

> A website to help with the export of routes from MapMyRun and the import to Garmin Connect.

## Status

Visit [mapmyrun-to-garmin.carlyzach.com](https://mapmyrun-to-garmin.carlyzach.com) to import MapMyRun routes to Garmin Connect!

As always, you can still use the [manual method](MANUAL.md).

Follow my repo for future updates!

## Developing

First start by downloading the code and opening that directory:
`git clone https://github.com/blaskovicz/mapmyrun-to-garmin.git`.

The app backend is written in [Golang](https://golang.org/) and can be started via `go run cmd/web/main.go`.
Set `PORT` in your environment to change the port it listens on. Set `ENVIRONMENT` to `development` to allow
CORS. Changes to this code require a manual restart.

The app frontend is written in [Vue.js](https://vuejs.org) and can be started and developed using the `yarn start` command.
To tell the frontend what backend to talk to, you can set `VUE_APP_MMR_API=http://some-url/` (this is only needed during development).
In production, the backend will serve the compiled frontend files (which can be tested by running `yarn build` and then visitting
the URI of the backend).

## Running on Heroku

1.  `$ heroku apps:create --buildpack heroku/go # create an app, set git remote`
2.  `$ heroku buildpacks:add heroku/nodejs # we need node as well`
3.  `$ git push heroku master # deploy our app`
4.  `$ curl -s -o /dev/null -w "%{http_code}" $(heroku info -j | jq -r '.app.web_url')/routes/new # check the site status, expecting 200`

## Running on Docker

1.  `$ docker pull blaskovicz/mapmyrun-to-garmin # download the docker image`
2.  `$ CONTAINER_ID=$(docker run --name mapmyrun-to-garmin --restart always -d -p 0:80/tcp -e PORT=80 blaskovicz/mapmyrun-to-garmin) # daemonize the container with docker, saving the container id`
3.  `$ CONTAINER_PORT=$(docker inspect $CONTAINER_ID | jq -r '.[0].NetworkSettings.Ports["80/tcp"][0].HostPort') # extract the external port mapping`
4.  `$ curl -s -o /dev/null -w "%{http_code}" 0.0.0.0:$CONTAINER_PORT/routes/new # check the site status, expecting 200`
