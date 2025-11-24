<script lang="ts">
  import { onMount } from "svelte";
  import { schemas } from "../stores/schemaStore";
  import { currentUser } from "../stores/userStore";
  import type { Schema } from "../lib/types";
  import VariableList from "../components/VariableList.svelte";

  export let user: string;
  export let params;
  let schema: Schema | null = null;

  onMount(async () => {
    try {
      await schemas.load(user);
      schema = $schemas.find((s) => s._id === params.id) || null;
    } catch (e) {
      const error = e instanceof Error ? e.message : "Failed to load schemas";
      console.error(error);
    }
  });

  async function saveSchema() {
    if (schema === null) {
      console.error("No schema to save");
      return;
    }
    await schemas.save(user, schema);
  }
</script>

{#if schema}
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
{/if}
