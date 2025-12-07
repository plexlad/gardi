<script lang="ts">
  import type { Instance, Schema, Visualization } from "../lib/types";
  
  export let instance: Instance;
  export let schema: Schema;
  
  // Get the visualization template from schema
  $: visualizations = schema.visualizations || [];
  
  function getVariableValue(varKey: string): any {
    // Return instance value if set, otherwise return schema default
    if (instance.variables && instance.variables[varKey] !== undefined) {
      return instance.variables[varKey];
    }
    const schemaVar = schema.variables[varKey];
    return schemaVar?.default ?? '';
  }
  
  function formatValue(varKey: string, value: any): string {
    const variable = schema.variables[varKey];
    if (!variable) return String(value);
    
    switch (variable.type) {
      case 'boolean':
        return value ? 'Yes' : 'No';
      case 'number':
        return value != null ? Number(value).toFixed(2) : 'N/A';
      case 'array':
        return Array.isArray(value) ? value.join(', ') : 'N/A';
      default:
        return String(value);
    }
  }
  
  function updateVariable(varKey: string, value: any) {
    if (!instance.variables) {
      instance.variables = {};
    }
    instance.variables[varKey] = value;
    // Trigger reactivity
    instance = instance;
  }
</script>

<div class="instance-display">
  {#if visualizations.length === 0}
    <div class="no-viz">
      <p>No visualization template defined for this schema.</p>
      <p class="hint">Edit the schema to create a visualization template.</p>
    </div>
  {:else}
    <div class="visualizations">
      {#each visualizations as viz}
        {#if viz.type === "grid"}
          <div class="grid-viz">
            <h3>{viz.name}</h3>
            <div
              class="grid"
              style="grid-template-columns: repeat({viz.config.cols}, 1fr);"
            >
              {#each viz.child_visualizations || [] as cell}
                <div class="grid-cell {cell.type}">
                  {#if cell.type === "empty"}
                    <span class="empty-label">—</span>
                  {:else if cell.type === "variable"}
                    {@const varKey = cell.config.variable_key}
                    {@const variable = schema.variables[varKey]}
                    {@const value = getVariableValue(varKey)}
                    
                    <div class="cell-content">
                      <div class="cell-label">{varKey}</div>
                      <div class="cell-input">
                        {#if variable.type === "number"}
                          <input
                            type="number"
                            bind:value={instance.variables[varKey]}
                            min={variable.min}
                            max={variable.max}
                            placeholder={variable.default}
                          />
                        {:else if variable.type === "boolean"}
                          <label class="checkbox-label">
                            <input
                              type="checkbox"
                              bind:checked={instance.variables[varKey]}
                            />
                            <span>{value ? 'Yes' : 'No'}</span>
                          </label>
                        {:else if variable.type === "enum"}
                          <select bind:value={instance.variables[varKey]}>
                            {#if !value}
                              <option value="">Select...</option>
                            {/if}
                            {#each variable.options || [] as opt}
                              <option value={opt}>{opt}</option>
                            {/each}
                          </select>
                        {:else}
                          <input
                            type="text"
                            bind:value={instance.variables[varKey]}
                            placeholder={variable.default}
                          />
                        {/if}
                      </div>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          </div>
        {/if}
      {/each}
    </div>
  {/if}
</div>

<style>
  .instance-display {
    width: 100%;
  }
  
  .no-viz {
    text-align: center;
    padding: 3rem 1rem;
    color: #666;
  }
  
  .no-viz p {
    margin: 0.5rem 0;
  }
  
  .hint {
    font-size: 0.9rem;
    color: #999;
  }
  
  .visualizations {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
  
  .grid-viz {
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
    padding: 1.5rem;
    background: #fafafa;
  }
  
  .grid-viz h3 {
    margin: 0 0 1.5rem 0;
    color: #333;
    font-size: 1.2rem;
  }
  
  .grid {
    display: grid;
    gap: 1rem;
  }
  
  .grid-cell {
    background: white;
    border: var(--border-width) solid var(--border-color);
    border-radius: 6px;
    min-height: 120px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }
  
  .grid-cell.empty {
    border-style: dashed;
    background: #f9f9f9;
  }
  
  .empty-label {
    color: #ccc;
    font-size: 1.5rem;
  }
  
  .grid-cell.variable {
    padding: 1rem;
    background: white;
  }
  
  .grid-cell.variable:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .cell-content {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .cell-label {
    font-weight: 600;
    color: #555;
    font-size: 0.9rem;
    text-transform: capitalize;
  }
  
  .cell-input {
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  input[type="text"],
  input[type="number"],
  select {
    width: 100%;
    padding: 0.75rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    font-size: 1rem;
    font-family: inherit;
    background: white;
    transition: border-color 0.2s;
  }
  
  input[type="text"]:focus,
  input[type="number"]:focus,
  select:focus {
    outline: none;
    border-color: #4a9eff;
    box-shadow: 0 0 0 3px rgba(74, 158, 255, 0.1);
  }
  
  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
    font-size: 1rem;
    color: #333;
  }
  
  input[type="checkbox"] {
    width: 24px;
    height: 24px;
    cursor: pointer;
    accent-color: #4a9eff;
  }
  
  .checkbox-label span {
    font-weight: 500;
  }
  
  select {
    cursor: pointer;
  }
</style>
