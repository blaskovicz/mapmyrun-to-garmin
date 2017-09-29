# Importing Routes from MapMyRun to Garmin Manually

Currently, routes can be exported from [MapMyRun](www.mapmyrun.com) in GPX format and then subsequently imported to [Garmin Connect](connect.garmin.com).

Note that you can now use mapmyrun-to-garmin.herokuapp.com (my website) to automatically convert routes.

If you still would like to do it manually, read on...

## Steps

1) Create the route on www.mapmyrun.com. Take note of the route ID in the url (eg: www.mapmyrun.com/routes/view/123 would have route ID `123`).
2) Note your MapMyRun.com `auth-token` cookie (eg: `US.76...`).
3) Visit [Garmin Connect](https://connect.garmin.com/modern/), log in, and note your connect.garmin.com `SESSION` (eg: `54fd-234...`) and `pin.m` (eg: `59db...`) cookies.
4) Download the latest [mapmyrun-to-garmin](https://github.com/blaskovicz/mapmyrun-to-garmin/tags) binary for your OS.
5) Run the tool in a command prompt, making sure to fill out the args with actual values:

```term
$ ./mapmyrun-to-garmin --route=$routeID --underarmour-cookie-auth-token=$underArmourAuthTokenCookie --garmin-session=$garminSessionCookie --garmin-pin-m=$garminPinMCookie
```

6) If this completes successfully, you will get a message with a link to the course in the connect dashboard. Upon failure, please ensure you have the most up-to-date cookies.

7) To push the course to your device (eg: Forerunner or Fenix), open the Garmin Connect course overview and click 'Send to Device.' On the PC this requires Garmin Express whereas on mobile, it will sync automatically over Bluetooth (recommended).

**Congratulations**, you've uploaded your course to Garmin and/or your device!
Please let me know if any corrections should be made to this article (pull requests and issues welcome).

My future plan includes automating the binary execution via a website, so stay tuned to this repo!
