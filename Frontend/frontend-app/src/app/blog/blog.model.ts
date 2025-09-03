import { Like } from "./like.model";

export interface Blog {
  id: string;
  userId: string;
  title: string;
  description: string;
  date_of_creation: string;
  images: string[];
  likes: Like[];           // array of Like objects
}
