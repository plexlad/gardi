<script lang="ts">
  import { VariableTypes } from "../lib/types";
  export let variable: any;
  export let remove: () => void;

  $: optionsString = variable.options ? variable.options.join(", ") : "";
</script>

<div class="variable-card">
  <select bind:value={variable.type}>
    {#each Object.values(VariableTypes) as t};
      <option value={t}>{t}</option>
    {/each}
  </select>
  <label>
    Default:
    <input bind:value={variable.default} />
  </label>
  {#if variable.type === VariableTypes.Number}
    <label>
      Min:
      <input type="number" bind:value={variable.min} />
    </label>
    <label>
      Max:
      <input type="number" bind:value={variable.max} />
    </label>
  {:else if variable.type === VariableTypes.Enum}
    <label>
      Options (comma separated):
      <input
        bind:value={optionsString}
        on:input={() => {
          variable.options = optionsString
            .split(",")
            .map((opt: string) => opt.trim())
            .filter((opt: string) => opt.length > 0);
        }}
      />
    </label>
  {:else if variable.type === VariableTypes.Array}
    <label>
      Items Type:
      <select bind:value={variable.items}>
        {#each Object.values(VariableTypes) as t};
          <option value={t}>{t}</option>
        {/each}
      </select>
    </label>
  {/if}
  <button on:click={remove}>Remove</button>
</div>
