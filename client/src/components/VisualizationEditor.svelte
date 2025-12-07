<script lang="ts">
  import { onMount } from "svelte";
  import type { Schema, Visualization, VisualizationTypes } from "../lib/types";
  
  export let schema: Schema;
  
  let visualizations: Visualization[] = [];
  let draggedItem: { type: string; data: any } | null = null;
  let editingGrid: { index: number; rows: number; cols: number } | null = null;
  
  // Available visualization types to drag from palette
  const paletteItems = [
    { type: "grid", label: "Grid" },
  ];
  
  // Get all variables from schema for the variable palette
  $: variableList = Object.entries(schema.variables || {}).map(([key, variable]) => ({
    key,
    variable
  }));
  
  function handlePaletteDragStart(e: DragEvent, item: any) {
    draggedItem = { type: "palette", data: item };
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = "copy";
    }
  }
  
  function handleVariableDragStart(e: DragEvent, varKey: string) {
    draggedItem = { type: "variable", data: varKey };
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = "copy";
    }
  }
  
  function handleCanvasDrop(e: DragEvent) {
    e.preventDefault();
    if (!draggedItem) return;
    
    if (draggedItem.type === "palette" && draggedItem.data.type === "grid") {
      // Add a new grid to the canvas
      const newGrid: Visualization = {
        name: `Grid ${visualizations.length + 1}`,
        type: "grid" as VisualizationTypes,
        child_visualizations: createEmptyGrid(3, 3),
        config: { rows: 3, cols: 3 }
      };
      visualizations = [...visualizations, newGrid];
    }
    
    draggedItem = null;
  }
  
  function handleGridCellDrop(e: DragEvent, gridIndex: number, cellIndex: number) {
    e.preventDefault();
    e.stopPropagation();
    
    if (!draggedItem || draggedItem.type !== "variable") return;
    
    const grid = visualizations[gridIndex];
    if (!grid.child_visualizations) return;
    
    // Replace the empty cell with a variable visualization
    grid.child_visualizations[cellIndex] = {
      name: draggedItem.data,
      type: "variable" as VisualizationTypes,
      config: { variable_key: draggedItem.data }
    };
    
    visualizations = [...visualizations];
    draggedItem = null;
  }
  
  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = "copy";
    }
  }
  
  function createEmptyGrid(rows: number, cols: number): Visualization[] {
    const cells: Visualization[] = [];
    for (let i = 0; i < rows * cols; i++) {
      cells.push({
        name: `Empty ${i}`,
        type: "empty" as VisualizationTypes,
        config: {}
      });
    }
    return cells;
  }
  
  function deleteVisualization(index: number) {
    visualizations = visualizations.filter((_, i) => i !== index);
  }
  
  function clearCell(gridIndex: number, cellIndex: number) {
    const grid = visualizations[gridIndex];
    if (!grid.child_visualizations) return;
    
    grid.child_visualizations[cellIndex] = {
      name: `Empty ${cellIndex}`,
      type: "empty" as VisualizationTypes,
      config: {}
    };
    
    visualizations = [...visualizations];
  }
  
  function openGridEditor(index: number) {
    const grid = visualizations[index];
    editingGrid = {
      index,
      rows: grid.config.rows,
      cols: grid.config.cols
    };
  }
  
  function applyGridResize() {
    if (!editingGrid) return;
    
    const grid = visualizations[editingGrid.index];
    const oldRows = grid.config.rows;
    const oldCols = grid.config.cols;
    const newRows = editingGrid.rows;
    const newCols = editingGrid.cols;
    
    // Create new grid with preserved cells where possible
    const newCells: Visualization[] = [];
    for (let r = 0; r < newRows; r++) {
      for (let c = 0; c < newCols; c++) {
        const oldIndex = r * oldCols + c;
        if (r < oldRows && c < oldCols && grid.child_visualizations && grid.child_visualizations[oldIndex]) {
          newCells.push(grid.child_visualizations[oldIndex]);
        } else {
          newCells.push({
            name: `Empty ${r * newCols + c}`,
            type: "empty" as VisualizationTypes,
            config: {}
          });
        }
      }
    }
    
    grid.config = { rows: newRows, cols: newCols };
    grid.child_visualizations = newCells;
    visualizations = [...visualizations];
    editingGrid = null;
  }
  
  function getVariableDisplay(varKey: string): string {
    // For schema editor, just show the variable type
    const variable = schema.variables[varKey];
    if (!variable) return "N/A";
    return variable.type;
  }
  
  // Export/save the visualization configuration
  export function getVisualizationConfig(): Visualization[] {
    return visualizations;
  }
  
  // Load existing visualization configuration
  export function loadVisualizationConfig(config: Visualization[]) {
    visualizations = config;
  }
</script>

<div class="editor-container">
  <!-- Left Sidebar: Palette -->
  <aside class="palette">
    <h3>Visualizations</h3>
    <div class="palette-items">
      {#each paletteItems as item}
        <div
          class="palette-item"
          draggable="true"
          on:dragstart={(e) => handlePaletteDragStart(e, item)}
        >
          <div class="grid-icon"></div>
          <span>{item.label}</span>
        </div>
      {/each}
    </div>
    
    <h3>Variables</h3>
    <div class="variable-list">
      {#each variableList as { key, variable }}
        <div
          class="variable-item"
          draggable="true"
          on:dragstart={(e) => handleVariableDragStart(e, key)}
          title={`${key} (${variable.type})`}
        >
          <div class="var-icon"></div>
          <span class="var-name">{key}</span>
          <span class="var-type">{variable.type}</span>
        </div>
      {/each}
    </div>
  </aside>
  
  <!-- Main Canvas -->
  <main
    class="canvas"
    on:drop={handleCanvasDrop}
    on:dragover={handleDragOver}
  >
    <div class="canvas-header">
      <h2>Visualization Template</h2>
      <p class="hint">Design how data will be displayed for instances of this schema</p>
    </div>
    
    {#if visualizations.length === 0}
      <div class="empty-canvas">
        <p>Drop a grid here to get started</p>
      </div>
    {/if}
    
    <div class="visualizations">
      {#each visualizations as viz, vizIndex}
        {#if viz.type === "grid"}
          <div class="grid-container">
            <div class="grid-header">
              <h4>{viz.name}</h4>
              <div class="grid-actions">
                <button
                  class="btn-edit"
                  on:click={() => openGridEditor(vizIndex)}
                  title="Edit grid size"
                >
                  Edit
                </button>
                <button
                  class="btn-delete"
                  on:click={() => deleteVisualization(vizIndex)}
                  title="Delete grid"
                >
                  Delete
                </button>
              </div>
            </div>
            <div
              class="grid"
              style="grid-template-columns: repeat({viz.config.cols}, 1fr);"
            >
              {#each viz.child_visualizations || [] as cell, cellIndex}
                <div
                  class="grid-cell {cell.type}"
                  on:drop={(e) => handleGridCellDrop(e, vizIndex, cellIndex)}
                  on:dragover={handleDragOver}
                >
                  {#if cell.type === "empty"}
                    <span class="empty-label">Empty</span>
                  {:else if cell.type === "variable"}
                    <div class="cell-content">
                      <div class="cell-header">
                        <span class="cell-var-name">{cell.config.variable_key}</span>
                        <button
                          class="btn-clear"
                          on:click={() => clearCell(vizIndex, cellIndex)}
                          title="Clear cell"
                        >
                          ×
                        </button>
                      </div>
                      <div class="cell-value">
                        {getVariableDisplay(cell.config.variable_key)}
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
  </main>
</div>

<!-- Grid Size Editor Modal -->
{#if editingGrid}
  <div class="modal-overlay" on:click={() => editingGrid = null}>
    <div class="modal" on:click|stopPropagation>
      <h3>Edit Grid Size</h3>
      <div class="form-group">
        <label>
          Rows:
          <input
            type="number"
            min="1"
            max="10"
            bind:value={editingGrid.rows}
          />
        </label>
        <label>
          Columns:
          <input
            type="number"
            min="1"
            max="10"
            bind:value={editingGrid.cols}
          />
        </label>
      </div>
      <div class="modal-actions">
        <button class="btn-cancel" on:click={() => editingGrid = null}>
          Cancel
        </button>
        <button class="btn-apply" on:click={applyGridResize}>
          Apply
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .editor-container {
    display: flex;
    height: calc(100vh - 60px);
    gap: 1rem;
    padding: 1rem;
  }
  
  .palette {
    width: 280px;
    background: #f8f9fa;
    border-radius: 8px;
    padding: 1rem;
    overflow-y: auto;
  }
  
  .palette h3 {
    margin: 0 0 1rem 0;
    color: #333;
    font-size: 1rem;
    font-weight: 600;
  }
  
  .palette-items {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 2rem;
  }
  
  .palette-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    background: white;
    border: 2px solid #e0e0e0;
    border-radius: 6px;
    cursor: move;
    transition: all 0.2s;
  }
  
  .palette-item:hover {
    border-color: #4a9eff;
    box-shadow: 0 2px 8px rgba(74, 158, 255, 0.2);
  }
  
  .grid-icon {
    width: 24px;
    height: 24px;
    background: var(--button);
    border-radius: 4px;
    position: relative;
  }
  
  .grid-icon::before,
  .grid-icon::after {
    content: '';
    position: absolute;
    background: white;
  }
  
  .grid-icon::before {
    width: 18px;
    height: 2px;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }
  
  .grid-icon::after {
    width: 2px;
    height: 18px;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }
  
  .icon {
    font-size: 1.5rem;
  }
  
  .variable-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .variable-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    background: white;
    border: 2px solid #e0e0e0;
    border-radius: 6px;
    cursor: move;
    transition: all 0.2s;
  }
  
  .variable-item:hover {
    border-color: #4caf50;
    box-shadow: 0 2px 8px rgba(76, 175, 80, 0.2);
  }
  
  .var-icon {
    width: 16px;
    height: 16px;
    background: #4caf50;
    border-radius: 50%;
    flex-shrink: 0;
  }
  
  .var-name {
    flex: 1;
    font-weight: 500;
    color: #333;
    font-size: 0.9rem;
  }
  
  .var-type {
    font-size: 0.75rem;
    color: #666;
    background: #f0f0f0;
    padding: 0.2rem 0.4rem;
    border-radius: 3px;
  }
  
  .canvas {
    flex: 1;
    background: white;
    border: 2px dashed #ddd;
    border-radius: 8px;
    padding: 1.5rem;
    overflow-y: auto;
  }
  
  .canvas-header {
    margin-bottom: 2rem;
  }
  
  .canvas-header h2 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }
  
  .hint {
    color: #666;
    font-size: 0.9rem;
    margin: 0;
  }
  
  .empty-canvas {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 300px;
    border: 2px dashed #ccc;
    border-radius: 8px;
    color: #999;
    font-size: 1.1rem;
  }
  
  .visualizations {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .grid-container {
    border: var(--border-width) solid var(--border-color);
    border-radius: 8px;
    padding: 1rem;
    background: #fafafa;
  }
  
  .grid-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .grid-header h4 {
    margin: 0;
    color: #333;
  }
  
  .grid-actions {
    display: flex;
    gap: 0.5rem;
  }
  
  .btn-edit, .btn-delete {
    background: white;
    border: var(--border-width) solid var(--border-color);
    font-size: 0.85rem;
    cursor: pointer;
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    transition: all 0.2s;
    color: #555;
  }
  
  .btn-edit:hover {
    background: #f0f0f0;
    border-color: #4a9eff;
    color: #4a9eff;
  }
  
  .btn-delete:hover {
    background: #fff0f0;
    border-color: #f44336;
    color: #f44336;
  }
  
  .grid {
    display: grid;
    gap: 0.5rem;
  }
  
  .grid-cell {
    aspect-ratio: 1;
    border: 2px solid #e0e0e0;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    transition: all 0.2s;
    min-height: 100px;
  }
  
  .grid-cell.empty {
    border-style: dashed;
    cursor: pointer;
  }
  
  .grid-cell.empty:hover {
    border-color: #4caf50;
    background: #f0f9f0;
  }
  
  .empty-label {
    color: #999;
    font-size: 0.9rem;
  }
  
  .grid-cell.variable {
    padding: 0.75rem;
    background: var(--button);
    border: none;
  }
  
  .cell-content {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  
  .cell-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 0.5rem;
  }
  
  .cell-var-name {
    color: white;
    font-weight: 600;
    font-size: 0.85rem;
    word-break: break-word;
  }
  
  .btn-clear {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: none;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    cursor: pointer;
    font-size: 1rem;
    line-height: 1;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .btn-clear:hover {
    background: rgba(255, 255, 255, 0.3);
  }
  
  .cell-value {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 1.5rem;
    font-weight: 700;
  }
  
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  
  .modal {
    background: white;
    border-radius: 8px;
    padding: 2rem;
    min-width: 400px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  }
  
  .modal h3 {
    margin: 0 0 1.5rem 0;
    color: #333;
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }
  
  .form-group label {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    font-weight: 500;
    color: #555;
  }
  
  .form-group input {
    padding: 0.5rem;
    border: var(--border-width) solid var(--border-color);
    border-radius: 4px;
    font-size: 1rem;
  }
  
  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }
  
  .btn-cancel, .btn-apply {
    padding: 0.5rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .btn-cancel {
    background: #e0e0e0;
    color: #333;
  }
  
  .btn-cancel:hover {
    background: #d0d0d0;
  }
  
  .btn-apply {
    background: #4a9eff;
    color: white;
  }
  
  .btn-apply:hover {
    background: #3a8eef;
  }
</style>
