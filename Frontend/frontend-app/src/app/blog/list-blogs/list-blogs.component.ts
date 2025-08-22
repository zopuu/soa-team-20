import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  ChangeDetectionStrategy,
} from '@angular/core';
import { BlogService } from '../blog.service';
import { Blog } from '../blog.model';

@Component({
  selector: 'app-blog-list',
  templateUrl: './list-blogs.component.html',
  styleUrls: ['./list-blogs.component.css'],
})
export class ListBlogsComponent implements OnChanges {
  @Input() userId?: string;

  blogs: Blog[] = [];
  loading = false;
  error = '';

  constructor(private blogService: BlogService) {}

  ngOnChanges(_: SimpleChanges): void {
    this.load();
  }

  private load(): void {
    this.loading = true;
    this.error = '';
    const req$ = this.userId
      ? this.blogService.getAllByUser(this.userId)
      : this.blogService.getAllBlogs();

    req$.subscribe({
      next: (data) => {
        this.blogs = data || [];
        console.log(this.blogs);
        this.loading = false;
      },
      error: (err) => {
        console.error('Failed to load blogs', err);
        this.error = 'Failed to load blogs.';
        this.loading = false;
      },
    });
  }

  trackById = (_: number, b: Blog) => b.id;
}
