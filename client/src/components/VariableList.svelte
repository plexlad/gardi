<script lang="ts">
  import type { Variable } from "../lib/types";
  import { VariableTypes } from "../lib/types";
  import VariableEditor from "./VariableEditor.svelte";
  
  export let variables: { [key: string]: Variable } = {};
  let newVarName: string = "";
  
  function addVar() {
    if (!newVarName || variables[newVarName]) return;
    variables = {
      ...variables,
      [newVarName]: {
        type: VariableTypes.String,
        default: "",
        min: null,
        max: null,
        options: null,
        items: null,
      },
    };
    newVarName = "";
  }
  
  function removeVar(key: string) {
    const { [key]: _, ...rest } = variables;
    variables = rest;
  }

  let filterText: string = "";
  $: filteredKeys = Object.keys(variables).filter(key =>
    key.toLowerCase().includes(filterText.toLowerCase())
  );
</script>

<div class="variables-list">
  <div class="filter-bar">
    <input 
      type="text" 
      placeholder="Filter variables..." 
      bind:value={filterText} 
      class="filter-input"
    />
  </div>

  <div class="add-variable">
    <input 
      type="text"
      placeholder="Variable name" 
      bind:value={newVarName}
      on:keydown={(e) => e.key === 'Enter' && addVar()}
    />
    <button on:click={addVar}>Add Variable</button>
  </div>

  {#if filteredKeys.length === 0}
    <p class="empty-state">No variables found.</p>
  {/if}
  
  {#if Object.keys(variables).length === 0}
    <p class="empty-state">No variables yet. Add one to get started!</p>
  {/if}
  
  {#each filteredKeys as key (key)}
    <div class="variable-item">
      <div class="variable-header">
        <h4>{key}</h4>
        <button class="delete-btn" on:click={() => removeVar(key)}>Delete</button>
      </div>
      <VariableEditor
        bind:variable={variables[key]} 
        remove={() => removeVar(key)}
      />
    </div>
  {/each}
</div>

<style>
  .variables-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .add-variable {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
  
  .add-variable input {
    flex: 1;
    padding: 0.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    font-size: 0.95rem;
  }
  
  .add-variable input:focus {
    outline: none;
    border-color: #4a9eff;
  }
  
  .add-variable button {
    padding: 0.5rem 1.5rem;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
    white-space: nowrap;
  }

  button {
    background-color: var(--button);
  }
  
  button:hover {
    background-color: var(--button-hover);
  }
  
  .variable-item {
    border: var(--border-width) solid var(--border-color);
    padding: 1rem;
    border-radius: 8px;
    background-color: #fafafa;
  }
  
  .variable-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  h4 {
    margin: 0;
    color: #333;
    font-size: 1.1rem;
  }
  
  .delete-btn {
    padding: 0.4rem 1rem;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background-color 0.2s;
  }
  
  .empty-state {
    text-align: center;
    color: #999;
    padding: 2rem;
    font-style: italic;
  }

  .filter-bar {
    margin-bottom: 0.5rem;
    display: flex;
  }
  .filter-input {
    flex: 1;
    padding: 0.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    font-size: 0.95rem;
  }
  .filter-input:focus {
    outline: none;
    border-color: #4a9eff;
  }
</style>
