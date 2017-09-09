# Importing Routes from MapMyRun to Garmin Manually

Currently, routes can be exported from [MapMyRun](www.mapmyrun.com) in GPX format and then subsequently imported to [Garmin Connect](connect.garmin.com).
Garmin Connect doesn't allow direct GPX route import, so we must first import it as an activity and then tell
Garmin to persist the route from the activity.

## Steps

1) Create the route on www.mapmyrun.com.
2) Under the 'Route Info' section, click 'Export this Route' and then 'Download GPX File'.
3) Open the GPX file in a text editor (I recommend [Notepad++](https://notepad-plus-plus.org/) or [Sublime Text](https://www.sublimetext.com/3)).
4) Change the opening `<gpx xmlns="http://www.topografix.com/GPX/1/1">` tag to:
```xml
<gpx creator="Garmin Connect" version="1.1"
  xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/11.xsd"
  xmlns="http://www.topografix.com/GPX/1/1"
  xmlns:ns3="http://www.garmin.com/xmlschemas/TrackPointExtension/v1"
  xmlns:ns2="http://www.garmin.com/xmlschemas/GpxExtensions/v3"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
```
5) After the newly changed `<gpx>` tag but before the `<trk>` tag, insert a `<metadata>` block of the form:
```xml
  <metadata>
    <link href="connect.garmin.com">
      <text>Garmin Connect</text>
    </link>
    <time>2017-09-08T22:20:18.000Z</time>
  </metadata>
```
6) In the `<trk>` section after the `</name>` tag, insert `<type>running</type>`.
7) Lastly, you will see various `<trkpt>` tags of the form `<trkpt lat="40.147" lon="-70.430"/>`. We must spoof
a time for all of these points. Luckily we can just use the same time as the start of the activity for all points.
In your edittor, do a find-and-replace of all text matching `"/>` with `"><time>2017-09-08T22:20:33.000Z</time></trkpt>`.
8) Save the file
9) Visit Garmin's [data-importer](https://connect.garmin.com/modern/import-data) after login to Garmin Connect.
10) Drag and drop the route file you editted and then click 'Import Data'. If an error occurs, ensure that you've followed
the steps above, otherwise, view the newly created "activity."
11) On the activity, click the gear icon and select 'Save as Course'; select your course type and then 'Continue'. Name the course
and then select 'Save New Course.' At this point, your new course exists in garmin connect.
12) For bookkeeping, delete the "activity" that was temporarily created by finding it onthe 'All Activities' page and then clicking
the 'garbage can' icon, followed by 'Delete.'
13) To push the course to your device (eg: Forerunner or Fenix), find the course overview and then click 'Send to Device.' On the
PC this requires Garmin Express whereas on mobile, it will sync automatically over Bluetooth (recommended).

Congratulations, you've uploaded your course to Garmin and/or your device!
Please let me know if any corrections should be made to this article (pull requests and issues welcome).

My future plan includes automating the course export/import via a website, so stay tuned to this repo!
