<script lang="ts">
  import { VariableTypes } from "../lib/types";
  export let variable: any;
  export const remove: () => void = () => {};
  
  $: optionsString = variable.options ? variable.options.join(", ") : "";
</script>

<div class="variable-editor">
  <div class="form-group">
    <label for="type-select" class="form-label">Type</label>
    <select id="type-select" bind:value={variable.type} class="form-select">
      {#each Object.values(VariableTypes) as t}
        <option value={t}>{t}</option>
      {/each}
    </select>
  </div>

  <div class="form-group">
    <label for="default-input" class="form-label">Default Value</label>
    <input id="default-input" bind:value={variable.default} class="form-input" placeholder="Enter default value" />
  </div>

  {#if variable.type === VariableTypes.Number}
    <div class="form-row">
      <div class="form-group">
        <label for="min-input" class="form-label">Min</label>
        <input id="min-input" type="number" bind:value={variable.min} class="form-input" placeholder="Min" />
      </div>
      <div class="form-group">
        <label for="max-input" class="form-label">Max</label>
        <input id="max-input" type="number" bind:value={variable.max} class="form-input" placeholder="Max" />
      </div>
    </div>
  {:else if variable.type === VariableTypes.Enum}
    <div class="form-group">
      <label for="options-input" class="form-label">Options (comma separated)</label>
      <input
        id="options-input"
        bind:value={optionsString}
        class="form-input"
        placeholder="option1, option2, option3"
        on:input={() => {
          variable.options = optionsString
            .split(",")
            .map((opt: string) => opt.trim())
            .filter((opt: string) => opt.length > 0);
        }}
      />
    </div>
  {:else if variable.type === VariableTypes.Array}
    <div class="form-group">
      <label for="items-select" class="form-label">Items Type</label>
      <select id="items-select" bind:value={variable.items} class="form-select">
        {#each Object.values(VariableTypes) as t}
          <option value={t}>{t}</option>
        {/each}
      </select>
    </div>
  {/if}
</div>

<style>
  .variable-editor {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .form-label {
    font-size: 0.9rem;
    font-weight: 500;
    color: #555;
  }

  .form-input,
  .form-select {
    padding: 0.6rem;
    border: var(--button-color) solid var(--button-width);
    border-radius: 4px;
    font-size: 0.95rem;
    background-color: white;
    transition: border-color 0.2s;
  }

  .form-input:focus,
  .form-select:focus {
    outline: none;
    border-color: #4a9eff;
  }

  .form-select {
    cursor: pointer;
  }

  .form-input::placeholder {
    color: #aaa;
  }
</style>
