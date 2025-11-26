import { writable } from "svelte/store";
import type { Writable } from "svelte/store";
import * as api from "../lib/api";
import type { Schema, NewSchemaRequest } from "../lib/types";

interface SchemaStore extends Writable<Schema[]> {
  load: (user: string) => Promise<void>;
  create: (user: string, data: NewSchemaRequest) => Promise<Schema>;
  save: (user: string, schema: Schema) => Promise<void>;
  delete: (user: string, schema: Schema) => Promise<void>;
}

function createSchemaStore(): SchemaStore {
  const { subscribe, set, update } = writable<Schema[]>([]);

  return {
    subscribe,
    set,
    update,

    load: async (user: string) => {
      try {
        const schemas = await api.getSchemas(user);
        console.log("API returned schemas:", schemas);
        set(schemas);
      } catch (error) {
        console.error("Failed to load schemas:", error);
        throw error;
      }
    },

    create: async (user: string, data: NewSchemaRequest) => {
      try {
        const newSchema = await api.createSchema(user, data);
        update((schemas) => [...schemas, newSchema]);
        return newSchema;
      } catch (error) {
        console.error("Failed to create schema:", error);
        throw error;
      }
    },

    save: async (user: string, schema: Schema) => {
      try {
        await api.saveSchema(user, schema);
        update((schemas) =>
          schemas.map((s) => (s._id === schema._id ? schema : s)),
        );
      } catch (e) {
        console.error("Failed to save schema:", e);
        throw e;
      }
    },

    delete: async (user: string, schema: Schema) => {
      try {
        await api.deleteSchema(user, schema);
        update((schemas) => schemas.filter((s) => s._id !== schema._id));
      } catch(e) {
        console.error("Failed to delete schema:", e);
        throw e;
      }
    },
  };
}

export const schemas = createSchemaStore();
