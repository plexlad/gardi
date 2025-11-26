<script lang="ts">
  import { onMount } from "svelte";
  import { schemas } from "../stores/schemaStore";
  import { currentUser } from "../stores/userStore";
  import type { Schema } from "../lib/types";
  import VariableList from "../components/VariableList.svelte";

  export let params;
  let user = $currentUser;

  let loading = true;

  // make sure that the schema is updated as well as variables
  $: schema = $schemas.find((s) => s._id === params.id);
  $: if (schema && schema.variables === null) {
    schema.variables = schema.variables || {};
  }

  onMount(async () => {
    if (!$currentUser) {
      console.error("No current user");
      loading = false;
      return;
    }

    try {
      console.log("Loading schemas for user:", currentUser);
      await schemas.load($currentUser);
      console.log("Schemas loaded, store value:", $schemas);
    } catch (e) {
      const error = e instanceof Error ? e.message : "Failed to load schemas";
      console.error(error);
    } finally {
      loading = false;
    }
  });

  async function saveSchema() {
    if (!schema) {
      console.error("No schema to save");
      return;
    }
    if (!user) {
      console.error("No current user");
      return;
    }
    await schemas.save(user, schema);
  }
</script>

{#if loading}
  <p>Loading...</p>
{:else if schema}
  <h2>Edit Schema: {schema.name}</h2>
  <label>
    Name:
    <input bind:value={schema.name} />
  </label>
  <label>
    Description:
    <textarea bind:value={schema.description}></textarea>
  </label>
  <h3>Variables</h3>
  <VariableList bind:variables={schema.variables} />
  <button on:click={saveSchema}>Save</button>
{:else if $currentUser}
  <p>Schema not found.</p>
{:else}
  <p>Please log in to view and edit schemas.</p>
{/if}
