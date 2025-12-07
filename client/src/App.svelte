<script>
  import Router from "svelte-spa-router";
  import routes from "./routes";
  import SiteWrapper from "./components/SiteWrapper.svelte";
</script>

<SiteWrapper>
<Router {routes} />
</SiteWrapper>
