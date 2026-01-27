export interface FieldError {
  field: string;
  message: string;
}

export interface ApiErrorBody {
  message: string;
  path: string;
  timestamp: string;
  details?: FieldError[];
}

export class ApiRequestError extends Error {
  status: number;
  path: string;
  timestamp: string;
  details?: FieldError[];

  constructor(status: number, body: ApiErrorBody) {
    super(body.message);
    this.name = "ApiRequestError";
    this.status = status;
    this.path = body.path;
    this.timestamp = body.timestamp;
    this.details = body.details;
  }
}
