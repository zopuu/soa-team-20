export interface BlogDto {
  userId: string;
  title: string;
  description: string;
  images?: File[]; // Optional array of files for image upload
}
