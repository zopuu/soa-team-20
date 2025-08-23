import { TourDifficulty } from './tour.model';

export interface TourDto {
  authorId: string;
  title: string;
  description: string;
  // backend expects numeric enum for difficulty
  difficulty: number;
  tags: string[];
}
