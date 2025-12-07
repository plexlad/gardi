<script lang="ts">
  import { onMount } from "svelte";
  import { schemas } from "../stores/schemaStore";
  import type { NewSchemaRequest, Schema } from "../lib/types";

  export let user: string;

  let loading = true;
  let error = "";
  let showForm = false;
  let newSchema: NewSchemaRequest = {
    name: "",
    description: "",
  };

  onMount(async () => {
    try {
      await schemas.load(user);
      console.log("Loaded schemas:", $schemas);
      loading = false;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load schemas";
      loading = false;
    }
  });

  async function handleCreate() {
    try {
      await schemas.create(user, newSchema);
      newSchema = { name: "", description: "" };
      showForm = false;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to create schema";
    }
  }

  async function handleDelete(schema: Schema) {
    try {
      await schemas.delete(user, schema);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to delete schema";
    }
  }
</script>

<div class="schema-list">
  <div class="header">
    <h3>Schemas</h3>
    <button on:click={() => (showForm = !showForm)}>
      {showForm ? "Cancel" : "+ New Schema"}
    </button>
  </div>

  {#if error}
    <div class="error">{error}</div>
  {/if}

  {#if showForm}
    <form on:submit|preventDefault={handleCreate} class="create-form">
      <input
        type="text"
        bind:value={newSchema.name}
        placeholder="Schema name"
        required
      />
      <textarea
        bind:value={newSchema.description}
        placeholder="Description"
        rows="3"
      ></textarea>
      <button type="reset">Reset</button>
      <button type="submit">Create</button>
    </form>
  {/if}

  {#if loading}
    <div class="loading">Loading schemas...</div>
  {:else if $schemas.length === 0}
    <div class="empty">No schemas yet. Create one to get started!</div>
  {:else}
    <div class="grid">
      {#each $schemas as schema}
        <div class="card">
          <a href={import.meta.env.BASE_URL + "#/schemas/" + schema._id}><h4>{schema.name}</h4></a>
          <p>{schema.description}</p>
          <button on:click={() => handleDelete(schema)}>Delete</button>
          <div class="meta">
            <small>Version: {schema.version}</small>
            <small
              >Created: {new Date(
                schema.created_at,
              ).toLocaleDateString()}</small
            >
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .schema-list {
    width: 100%;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  h3 {
    margin: 0;
  }

  button {
    padding: 0.5rem 1rem;
    background-color: var(--button);
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
  }

  button:hover {
    background-color: var(--button-hover);
  }

  .create-form {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    padding: 1rem;
    background: #f9f9f9;
    border-radius: var(--border-radius);
    margin-bottom: 1rem;
  }

  input,
  textarea {
    padding: 0.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    font-size: 0.9rem;
    font-family: inherit;
  }

  .error {
    padding: 0.75rem;
    background: #ffebee;
    color: #c62828;
    border-radius: var(--border-radius);
    margin-bottom: 1rem;
  }

  .loading,
  .empty {
    text-align: center;
    padding: 2rem;
    color: #666;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1rem;
  }

  .card {
    padding: 1rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    background: white;
    transition: box-shadow 0.2s;
  }

  .card:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .card h4 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }

  .card p {
    margin: 0 0 1rem 0;
    color: #666;
    font-size: 0.9rem;
  }

  .card button {
    margin-bottom: 1em;
  }

  .meta {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    color: #999;
  }
</style>
