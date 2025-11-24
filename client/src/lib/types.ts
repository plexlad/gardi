export interface Schema {
  _id: string;
  version: number;
  user_version: number;
  name: string;
  description: string;
  variables: { [key: string]: Variable };
  created_at: string;
  updated_at: string;
}

export interface Instance {
  _id: string,
  schema_id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

// The types of selectable variables
export enum VariableTypes {
  Number = "number",
  String = "string",
  Boolean = "boolean",
  Enum = "enum",
  Array = "array",
}

// Determines what your variable is
export interface Variable {
  type: VariableTypes;
  default: any;
  min: number | null;
  max: number | null;
  options: string[] | null; // for enum
  items: Variable | null; // for array type
}

export interface NewSchemaRequest {
  name: string;
  description: string;
}

export interface NewInstanceRequest {
  name: string;
  description: string;
  schema_id: string;
}

export interface ApiError {
  error: string;
}
