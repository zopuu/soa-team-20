import { BlogImage } from './blog-image.model';

export interface Blog {
  id: string;
  userId: string;
  title: string;
  description: string;
  date_of_creation: string;
  images: BlogImage[];
  //likes: Like[];           // array of Like objects
}
