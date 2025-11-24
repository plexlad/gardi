import Home from './pages/Home.svelte';
// import SchemaList from './pages/SchemaList.svelte';
import SchemaEditor from './pages/SchemaEditor.svelte';

export default {
  "/": Home,
  "/schemas/:id": SchemaEditor,
}
