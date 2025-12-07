<!-- src/App.svelte -->
<script lang="ts">
  import { currentUser } from '../stores/userStore';
  import UserLogin from '../components/UserLogin.svelte';
  import SchemaList from '../components/SchemaList.svelte';
  import InstanceList from '../components/InstanceList.svelte';
</script>

<main>
  <header>
    <UserLogin />
  </header>

  {#if $currentUser}
    <div class="content">
      <section>
        <SchemaList user={$currentUser} />
      </section>
      
      <section>
        <InstanceList user={$currentUser} />
      </section>
    </div>
  {/if}
</main>

<style>
  main {
    max-width: var(--site-width);
    margin: 0 auto;
    padding: 1rem;
  }
  
  header {
    margin-bottom: 2rem;
  }
  
  .content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }
  
  section {
    padding: 1rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
  }
 </style>
