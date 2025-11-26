import type {
  Schema,
  Instance,
  NewSchemaRequest,
  NewInstanceRequest,
  ApiError,
} from "./types";

const API_BASE = "http://localhost:5499";

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error: ApiError = await response.json();
    throw new Error(error.error || "API request failed");
  }
  return response.json();
}

export async function getSchemaIDs(user: string): Promise<string[]> {
  const response = await fetch(`${API_BASE}/${user}/schemas`);
  const data = await handleResponse<string[]>(response);
  return data ?? [];
}

export async function getSchemas(user: string): Promise<Schema[]> {
  const ids = await getSchemaIDs(user);
  console.log("Schemas - ids:", ids, "type:", typeof ids, "isArray:", Array.isArray(ids));
  const schemas = await Promise.all(
    ids.map(async (id: string) => {
      return handleResponse<Schema>(
        await fetch(`${API_BASE}/${user}/schemas/${id}`),
      );
    }),
  );
  return schemas;
}

export async function createSchema(
  user: string,
  data: NewSchemaRequest,
): Promise<Schema> {
  const response = await fetch(`${API_BASE}/${user}/schemas/new`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  return handleResponse<Schema>(response);
}

export async function saveSchema(
  user: string,
  schema: Schema,
): Promise<string> {
  const response = await fetch(`${API_BASE}/${user}/schemas/save`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(schema),
  });
  return response.text();
}

export async function deleteSchema(
  user: string,
  schema: Schema,
): Promise<void> {
  const schemaDeleteLink = `${API_BASE}/${user}/schemas/${schema._id}`;
  console.log("Deleting schema:", {
    user,
    schemaId: schema._id,
    url: schemaDeleteLink,
  });
  const response = await fetch(schemaDeleteLink, {
    method: "DELETE",
  });

  if (!response.ok) {
    const text = await response.text();
    console.error("Delete schema response text:", response.status, text);
    throw new Error(`Failed to delete schema: ${response.status}`);
  }
}

// Instance API
export async function getInstanceIDs(user: string): Promise<string[]> {
  const response = await fetch(`${API_BASE}/${user}/instances`);
  const data = await handleResponse<string[]>(response);
  return data ?? [];
}

export async function getInstances(user: string): Promise<Instance[]> {
  const ids = await getInstanceIDs(user);
  console.log("Instances - ids:", ids, "type:", typeof ids, "isArray:", Array.isArray(ids));
  const instances = await Promise.all(
    ids.map(async (id: string) => {
      return handleResponse<Instance>(
        await fetch(`${API_BASE}/${user}/instances/${id}`),
      );
    }),
  );
  return instances;
}

export async function createInstance(
  user: string,
  data: NewInstanceRequest,
): Promise<Instance> {
  const response = await fetch(`${API_BASE}/${user}/instances/new`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  return handleResponse<Instance>(response);
}

export async function saveInstance(
  user: string,
  instance: Instance,
): Promise<string> {
  const response = await fetch(`${API_BASE}/${user}/instances/save`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(instance),
  });
  return response.text();
}

export async function deleteInstance(
  user: string,
  instance: Instance,
): Promise<void> {
  const response = await fetch(
    `${API_BASE}/${user}/instances/${instance._id}`,
    {
      method: "DELETE",
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete instance");
  }
}
