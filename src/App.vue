<template>
  <div id="app">
    <ConvertHomepage />
    <footer class="footer">
      <b-container>
        <b-row>
          <b-col>
            <p class="text-muted">Questions or Concerns? Check out my <a rel="noopener noreferrer" href="https://github.com/blaskovicz/mapmyrun-to-garmin" target="_blank">Github</a> repo.
            </p>
          </b-col>
        </b-row>
      </b-container>
    </footer>
  </div>
       
</template>

<script>
import Vue from "vue";
import BootstrapVue from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";

Vue.use(BootstrapVue);

import ConvertHomepage from "./components/ConvertHomepage.vue";

export default {
  name: "app",
  components: {
    ConvertHomepage
  }
};
</script>

<style>
footer {
  margin-top: 10px;
}
</style>
