<script lang="ts">
  import { onMount } from "svelte";
  import { instances } from "../stores/instanceStore";
  import { schemas } from "../stores/schemaStore";
  import { currentUser } from "../stores/userStore";
  import InstanceDisplay from "../components/InstanceDisplay.svelte";
  import type { Instance, Schema } from "../lib/types";

  export let params;

  let loading = true;
  let error = "";
  let instance: Instance | null = null;
  let schema: Schema | null = null;
  let user = $currentUser;
  let showBasicDetails = false;

  onMount(async () => {
    loading = true;
    error = "";
    if (!user) {
      error = "No user logged in.";
      loading = false;
      return;
    }
    try {
      await schemas.load(user);
      await instances.load(user);
      instance = $instances.find(i => i._id === params.id) || null;
      
      if (!instance) {
        error = "Instance not found.";
        loading = false;
        return;
      }
      
      schema = $schemas.find(s => s._id === instance.schema_id) || null;
      
      if (!schema) {
        error = "Schema not found.";
        loading = false;
        return;
      }
      
      // Initialize variables object if it doesn't exist
      if (!instance.variables) {
        instance.variables = {};
      }
      
      // Initialize all variables with defaults if not set
      Object.entries(schema.variables || {}).forEach(([key, variable]) => {
        if (instance.variables[key] === undefined || instance.variables[key] === null) {
          // Set default value based on type
          if (variable.default !== undefined && variable.default !== null) {
            instance.variables[key] = variable.default;
          } else {
            // Set sensible defaults based on type
            switch (variable.type) {
              case 'number':
                instance.variables[key] = variable.min ?? 0;
                break;
              case 'boolean':
                instance.variables[key] = false;
                break;
              case 'enum':
                instance.variables[key] = variable.options?.[0] ?? '';
                break;
              case 'array':
                instance.variables[key] = [];
                break;
              default:
                instance.variables[key] = '';
            }
          }
        }
      });
      
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load data";
    } finally {
      loading = false;
    }
  });

  async function saveInstance() {
    if (!instance || !user) return;
    try {
      await instances.save(user, instance);
      error = "";
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to save instance";
    }
  }
</script>

{#if loading}
  <div class="container"><p class="message">Loading...</p></div>
{:else if error}
  <div class="container"><p class="error">{error}</p></div>
{:else if instance && schema}
  <div class="container">
    <header>
      <div class="header-content">
        <div>
          <h2>{instance.name}</h2>
          <p class="schema-info">Schema: {schema.name}</p>
        </div>
        <button
          class="toggle-details"
          on:click={() => showBasicDetails = !showBasicDetails}
        >
          {showBasicDetails ? 'Hide' : 'Show'} Details
        </button>
      </div>
    </header>
    
    {#if showBasicDetails}
      <section class="basic-details">
        <h3>Basic Information</h3>
        <label>
          <span class="label-text">Name:</span>
          <input type="text" bind:value={instance.name} />
        </label>
        <label>
          <span class="label-text">Description:</span>
          <textarea bind:value={instance.description}></textarea>
        </label>
      </section>
    {/if}
    
    <section class="instance-data">
      <h3>Instance Data</h3>
      <InstanceDisplay bind:instance {schema} />
    </section>
  </div>
  
  <!-- Floating save button -->
  <button class="floating-save-btn" on:click={saveInstance}>
    Save Instance
  </button>
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
  
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
  }
  
  h2 {
    color: #333;
    margin: 0 0 0.5rem 0;
    font-size: 1.8rem;
  }
  
  h3 {
    color: #333;
    margin: 0 0 1rem 0;
    font-size: 1.2rem;
  }
  
  .schema-info {
    color: #666;
    margin: 0;
    font-size: 0.95rem;
  }
  
  .toggle-details {
    padding: 0.5rem 1rem;
    background: white;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    color: #555;
    transition: all 0.2s;
  }
  
  .toggle-details:hover {
    background: #f5f5f5;
    border-color: #4a9eff;
    color: #4a9eff;
  }
  
  .basic-details {
    padding: 1.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
    margin-bottom: 2rem;
    background: white;
  }
  
  .instance-data {
    padding: 1.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
    margin-bottom: 2rem;
    background: white;
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
  
  input[type="text"], textarea {
    padding: 0.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    font-size: 1rem;
    font-family: inherit;
  }
  
  input[type="text"]:focus, textarea:focus {
    outline: none;
    border-color: #4a9eff;
  }
  
  textarea {
    min-height: 80px;
    resize: vertical;
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
  
  .error {
    color: #c62828;
    background: #ffebee;
    padding: 1em;
    border-radius: 8px;
    margin-bottom: 1em;
    text-align: center;
  }
  
  .message {
    text-align: center;
    padding: 2rem;
    color: #666;
    font-size: 1.1rem;
  }
</style>
