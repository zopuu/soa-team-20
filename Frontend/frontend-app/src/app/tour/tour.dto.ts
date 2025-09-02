import { TourDifficulty } from './tour.model';

export interface TourDto {
  authorId: string;
  title: string;
  description: string;
  difficulty: number;
  tags: string[];
}
