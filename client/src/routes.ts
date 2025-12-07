import Home from './pages/Home.svelte';
// import SchemaList from './pages/SchemaList.svelte';
import SchemaEditor from './pages/SchemaEditor.svelte';
import InstanceEditor from './pages/InstanceEditor.svelte';
import Information from './pages/Information.svelte';

export default {
  "/": Information,
  "/stuff": Home,
  "/schemas/:id": SchemaEditor,
  "/instances/:id": InstanceEditor,
}
