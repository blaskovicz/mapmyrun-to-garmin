<template>
  <div id="convert">
    <b-jumbotron>
        <template slot="header">
          <a href="/" style="color:inherit">MapMyRun to Garmin</a>
        </template>
        <template slot="lead">
          Import <a href="http://mapmyrun.com" target="_blank" rel="noopener noreferrer">MapMyRun</a> Routes to <a href="http://connect.garmin.com" target="_blank" rel="noopener noreferrer">Garmin Connect</a>
        </template>
    </b-jumbotron>

    <b-container>
      <b-alert variant="danger" v-if="response && response.error" show dismissible>
        {{response.error}}
      </b-alert>
      <b-alert variant="success" v-else-if="response && response.success" show dismissible>
        {{response.success}}
      </b-alert>

      <b-form @submit.prevent="onSubmit" @reset.prevent="onReset">
        <b-form-group label-for="garminSession" :state="garminSessionState" :invalid-feedback="garminSessionInvalid" :valid-feedback="garminSessionValid">
          <template slot="label" >
            Garmin <code>SESSIONID</code> Cookie
          </template>
          <template slot="description">
            Find this on <a rel="noopener noreferrer" target="_blank" href="http://connect.garmin.com">Garmin Connect</a> under domain <i>connect.garmin.com</i> (eg: <code>59123-...</code>).
          </template>
          <b-form-input required :state="garminSessionState" type="text" name="garminSession" id="garminSession" v-model="garminSession"></b-form-input>
        </b-form-group>

        <b-form-group label-for="mapmyrunAuthToken" :state="mapmyrunAuthTokenState" :invalid-feedback="mapmyrunAuthTokenInvalid" :valid-feedback="mapmyrunAuthTokenValid">
          <template slot="label" >
            MapMyRun <code>auth-token</code> Cookie
          </template>
          <template slot="description">
            Find this on <a rel="noopener noreferrer" target="_blank" href="http://mapmyrun.com">MapMyRun</a> under domain <i>mapmyrun.com</i> (eg: <code>US.123...</code>).
          </template>
          <b-form-input required :state="mapmyrunAuthTokenState" type="text" name="mapmyrunAuthToken" id="mapmyrunAuthToken" v-model="mapmyrunAuthToken"></b-form-input>
        </b-form-group>

        <b-form-group label-for="mapmyrunRouteId" :state="mapmyrunRouteIdState" :invalid-feedback="mapmyrunRouteIdInvalid" :valid-feedback="mapmyrunRouteIdValid">					
          <template slot="label">
						MapMyRun Route ID
          </template>
          <template slot="description">
            Find this on <a rel="noopener noreferrer" target="_blank" href="https://mapmyrun.com/routes/my_routes">MapMyRun</a> (eg: http://mapmyrun.com/routes/view/123 would have ID <code>123</code>).
          </template>
          <b-form-input required :state="mapmyrunRouteIdState" type="text" name="mapmyrunRouteId" id="mapmyrunRouteId" v-model="mapmyrunRouteId"></b-form-input>
        </b-form-group>

        <b-button id='convert-button' :disabled="(response && response.pending) || !garminSessionState || !mapmyrunAuthTokenState || !mapmyrunRouteIdState" type="submit" variant="primary">Submit</b-button>
        <b-button :disabled="response && response.pending" type="reset" variant="default">Reset</b-button>
      </b-form>
    </b-container>
  </div>
</template>

<script>
import axios from "axios";
import urljoin from "url-join";
export default {
  name: "ConvertHomepage",
  data() {
    return {
      garminSession: "",
      mapmyrunAuthToken: "",
      mapmyrunRouteId: "",
      response: null
    };
  },
  created() {
    this.genericInvalid = `It looks like this isn't field is invalid.`;
  },
  methods: {
    onReset() {
      this.garminSession = "";
      this.mapmyrunAuthToken = "";
      this.mapmyrunRouteId = "";
      this.response = null;
    },
    async onSubmit() {
      try {
        this.response = { pending: true };
        const payload = {
          garminSession: this.garminSession,
          mapmyrunAuthToken: this.mapmyrunAuthToken,
          mapmyrunRouteId: +this.mapmyrunRouteId
        };
        const importUrl = urljoin(
          process.env.VUE_APP_MMR_API || "",
          "/api/garmin/import"
        );
        const response = await axios.post(importUrl, payload);
        const { courseName, courseId } = response.data;
        this.response = {
          success: `Successfully imported course https://connect.garmin.com/modern/course/${courseId} '${courseName}' into Garmin Connect from http://mapmyrun.com/routes/view/${
            payload.mapmyrunRouteId
          }. Don't forget to send the course to your device via Garmin Connect!`
        };
        this.mapmyrunRouteId = "";
      } catch (err) {
        console.warn(err);
        let error;
        if (err.response) {
          error = err.response.data.message;
        } else {
          error = err.toString();
        }
        this.response = { error };
      }
    }
  },
  computed: {
    garminSessionState() {
      if (this.garminSession === "") return null;
      return !(
        this.garminSession.length < 5 || !this.garminSession.includes("-")
      );
    },
    garminSessionInvalid() {
      return `${this.genericInvalid}; must be non-empty and include a '-'.`;
    },
    garminSessionValid() {
      return "";
    },
    mapmyrunAuthTokenState() {
      if (this.mapmyrunAuthToken === "") return null;
      return !(
        this.mapmyrunAuthToken.length < 5 ||
        !this.mapmyrunAuthToken.includes(".")
      );
    },
    mapmyrunAuthTokenInvalid() {
      return `${this.genericInvalid}; must be non-empty and include a '.'.`;
    },
    mapmyrunAuthTokenValid() {
      return "";
    },
    mapmyrunRouteIdState() {
      if (this.mapmyrunRouteId === "") return null;
      return /^\d+$/.test(this.mapmyrunRouteId);
    },
    mapmyrunRouteIdInvalid() {
      return `${this.genericInvalid}; must be a number.`;
    },
    mapmyrunRouteIdValid() {
      return "";
    }
  }
};
</script>

<style>
label.bold {
  font-weight: bold;
}
form .form-group label {
  font-weight: bold;
}
#convert-button {
  margin-right: 4px;
}
</style>
