export type ActionType = 
  | "critical"
  | "error"
  | "request"
  | "response"
  | "query";

export type InputType =
  | "text"
  | "password"
  | "select"
  | "checkbox"
  | "radio"
  | "dropzone"
  | "hidden"
  | "file"
  | "tel"
  | "email"
  | "url"
  | "number"
  | "range"
  | "date"
  | "time"
  | "color"
  | "combo"
  | "textarea";

export interface Input {
  name: string;
  type: InputType;
  placeholder: string;
  value?: string;
  options?: string[];
  required: boolean;
  disabled: boolean;
  readOnly: boolean;
  order: number;
}

export interface ProgressMessage {
  type: "progress";
  progress: number;
  message: string;
}

export interface ActionMessage {
  type: ActionType;
  data?: any;
  inputs: Input[];
  title: string;
  message: string;
}

export type WebSocketMessage = ProgressMessage | ActionMessage;
