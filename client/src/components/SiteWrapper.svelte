<script>
  import logo from "/logo-from-svg-repo.svg";
</script>

<nav>
  <a href={import.meta.env.BASE_URL}>
    <img
      class="logo"
      src={logo}
      alt="a set of circles and lines depicting nodes connecting to each other"
    />
  </a>
  <ul class="links">
    <li><a href={import.meta.env.BASE_URL}>home</a></li>
    <li><a href={import.meta.env.BASE_URL + "/#/stuff"}>your stuff</a></li>
  </ul>
</nav>
<slot />

<style>
  nav {
    position: sticky;
    display: flex;
    border: var(--border-color) 0.1em solid;
    border-radius: 1em;
    align-items: center;
    justify-content: space-between;
    padding: 0.3em;
    padding-right: 0.6em;
    max-width: var(--site-width);
    margin: auto;
  }

  .logo {
    height: 3em;
    transition: 0.2s;
  }

  .logo:hover {
    transform: rotate(20deg);
  }

  .links {
    list-style-type: none;
  }

  .links li {
    display: inline;
  }

  .links a {
    padding: 0.5em;
    border: var(--border-color) 0.1em solid;
    border-radius: calc(var(--border-radius) * 0.5);
    text-decoration: none;
    color: inherit;
    transition: 0.2s;
  }

  .links a:hover {
    background-color: var(--button);
    color: white;
    border-color: var(--button);
  }
</style>
