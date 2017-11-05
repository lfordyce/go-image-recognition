export interface Label {
  label: string;
  probability: number;
}

export interface Post {
  filename: string;
  labels: Array<Label>;
}
