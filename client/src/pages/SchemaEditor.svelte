<script lang="ts">
  import { onMount } from "svelte";
  import { schemas } from "../stores/schemaStore";
  import { currentUser } from "../stores/userStore";
  import VariableList from "../components/VariableList.svelte";
  import VisualizationEditor from "../components/VisualizationEditor.svelte";
  
  export let params;
  
  let user = $currentUser;
  let loading = true;
  let vizEditor: VisualizationEditor;
  let showVisualizationEditor = false;
  
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
  
  // Load visualizations when editor becomes available or schema changes
  $: if (vizEditor && schema?.visualizations) {
    vizEditor.loadVisualizationConfig(schema.visualizations);
  }
  
  async function saveSchema() {
    if (!schema) {
      console.error("No schema to save");
      return;
    }
    if (!user) {
      console.error("No current user");
      return;
    }
    
    // Get visualization config from editor
    if (vizEditor) {
      const vizConfig = vizEditor.getVisualizationConfig();
      console.log("Saving visualization config:", vizConfig);
      schema.visualizations = vizConfig;
    }
    
    await schemas.save(user, schema);
  }
</script>

{#if loading}
  <div class="container">
    <p class="message">Loading...</p>
  </div>
{:else if schema}
  <div class="container">
    <header>
      <h2>Edit Schema: {schema.name}</h2>
      <div class="tabs">
        <button 
          class="tab"
          class:active={!showVisualizationEditor}
          on:click={() => showVisualizationEditor = false}
        >
          Schema Details
        </button>
        <button 
          class="tab"
          class:active={showVisualizationEditor}
          on:click={() => showVisualizationEditor = true}
        >
          Visualization Template
        </button>
      </div>
    </header>
    
    {#if !showVisualizationEditor}
      <section class="schema-details">
        <label>
          <span class="label-text">Name:</span>
          <input type="text" bind:value={schema.name} />
        </label>
        <label>
          <span class="label-text">Description:</span>
          <textarea bind:value={schema.description}></textarea>
        </label>
      </section>
      
      <section class="variables-section">
        <h3>Variables</h3>
        <VariableList bind:variables={schema.variables} />
      </section>
    {:else}
      <section class="visualization-section">
        <VisualizationEditor bind:this={vizEditor} {schema} />
      </section>
    {/if}
  </div>
  
  <!-- Floating save button -->
  <button class="floating-save-btn" on:click={saveSchema}>
    Save Schema
  </button>
{:else if $currentUser}
  <div class="container">
    <p class="message">Schema not found.</p>
  </div>
{:else}
  <div class="container">
    <p class="message">Please log in to view and edit schemas.</p>
  </div>
{/if}

<style>
  .container {
    max-width: var(--site-width);
    margin: 0 auto;
    padding: 1rem;
    padding-bottom: 5rem;
  }
  
  header {
    margin-bottom: 2rem;
  }
  
  h2 {
    color: #333;
    margin: 0 0 1rem 0;
  }
  
  h3 {
    color: #333;
    margin: 0 0 1rem 0;
  }
  
  .tabs {
    display: flex;
    gap: 0.5rem;
    border-bottom: var(--border-width) solid var(--border-color);
  }
  
  .tab {
    padding: 0.75rem 1.5rem;
    background: none;
    border: none;
    border-bottom: 3px solid transparent;
    cursor: pointer;
    font-size: 1rem;
    color: #666;
    transition: all 0.2s;
    margin-bottom: -2px;
  }
  
  .tab:hover {
    color: #333;
    background: #f5f5f5;
  }
  
  .tab.active {
    color: #4a9eff;
    border-bottom-color: #4a9eff;
    font-weight: 600;
  }
  
  .schema-details {
    padding: 1.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
    margin-bottom: 2rem;
  }
  
  .variables-section {
    padding: 1.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
    margin-bottom: 2rem;
  }
  
  .visualization-section {
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
    overflow: hidden;
    margin-bottom: 2rem;
  }
  
  label {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  
  .label-text {
    font-weight: 500;
    color: #555;
  }
  
  input[type="text"] {
    padding: 0.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    font-size: 1rem;
  }
  
  input[type="text"]:focus {
    outline: none;
    border-color: #4a9eff;
  }
  
  textarea {
    padding: 0.5rem;
    border-radius: 4px;
    border: var(--border-width) solid var(--border-color);
    font-size: 1rem;
    min-height: 100px;
    resize: vertical;
    font-family: inherit;
  }
  
  textarea:focus {
    outline: none;
    border-color: #4a9eff;
  }
  
  .floating-save-btn {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    padding: 1rem 2rem;
    background-color: #4a9eff;
    color: white;
    border: none;
    border-radius: 50px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(74, 158, 255, 0.4);
    transition: all 0.3s ease;
    z-index: 1000;
  }
  
  .floating-save-btn:hover {
    background-color: #3a8eef;
    box-shadow: 0 6px 16px rgba(74, 158, 255, 0.5);
    transform: translateY(-2px);
  }
  
  .floating-save-btn:active {
    transform: translateY(0);
  }
  
  .message {
    text-align: center;
    padding: 2rem;
    color: #666;
    font-size: 1.1rem;
  }
</style>
