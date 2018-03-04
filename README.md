# mapmyrun-to-garmin
>A website to help with the export of routes from MapMyRun and the import to Garmin Connect.

## Status

Visit [mapmyrun-to-garmin.carlyzach.com](https://mapmyrun-to-garmin.carlyzach.com) to import MapMyRun routes to Garmin Connect!

As always, you can still use the [manual method](MANUAL.md).

Follow my repo for future updates!

## Developing

First start by downloading the codei and switching to that directory. Run `go get github.com/blaskovicz/mapmyrun-to-garmin/...` or simply `git clone https://github.com/blaskovicz/mapmyrun-to-garmin.git`. 

## Running on Heroku

1) `$ heroku apps:create --buildpack heroku/go # create an app, set git remotee`
2) `$ heroku config:set ENVIRONMENT=production COOKIE_KEY=$(uuidgen) CSRF_KEY=$(uuidgen) # set needed config on heroku`
3) `$ git push heroku master # deploy our app`
4) `$ curl -s -o /dev/null -w "%{http_code}" $(heroku info -j | jq -r '.app.web_url')/routes/new # check the site status, expecting 200`

## Running on Docker

1) `$ docker pull blaskovicz/mapmyrun-to-garmin # download the docker image`
2) `$ CONTAINER_ID=$(docker run --name mapmyrun-to-garmin --restart always -d -e COOKIE_KEY=$(uuidgen) -e CSRF_KEY=$(uuidgen) -e ENVIRONMENT=production -p 0:80/tcp) # daemonize the container with docker, saving the container id`
3) `$ CONTAINER_PORT=$(docker inspect $CONTAINER_ID | jq -r '.[0].NetworkSettings.Ports["80/tcp"][0].HostPort') # extract the external port mapping`
4) `$ curl -s -o /dev/null -w "%{http_code}" 0.0.0.0:$CONTAINER_PORT/routes/new # check the site status, expecting 200`
