export enum Priority {
  LOW = 0,
  MEDIUM = 1,
  HIGH = 2,
  URGENT = 3,
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface PaginatedTodos {
  data: Todo[];
  meta: PaginationMeta;
}

export interface Todo {
  id: string;
  title: string;
  description: string;
  priority: Priority;
  completed: boolean;

  user_id: string;

  created_at: string;
  updated_at: string;
}

export interface TodoRequest {
  title?: string;
  description?: string;
  priority?: Priority;
  completed?: boolean;
}
