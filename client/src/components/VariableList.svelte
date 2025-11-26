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
</script>

<div class="variables-list">
  <div class="add-variable">
    <input placeholder="Variable name" bind:value={newVarName} />
    <button on:click={addVar}>Add Variable</button>
  </div>
  {#each Object.entries(variables) as [key, _] (key)};
    <div class="variable-item">
      <h4>{key}</h4>
      <VariableEditor
        bind:variable={variables[key]} 
        remove={() => removeVar(key)}
      />
      <button on:click={() => removeVar(key)}>Delete</button>
    </div>
  {/each}
</div>

<style>
  .variables-list { display: flex; flex-direction: column; gap: 1rem; }
  .add-variable { display: flex; gap: 0.5rem; }
  .variable-item { border: 1px solid #ddd; padding: 0.5rem; border-radius: 4px; }
</style>
