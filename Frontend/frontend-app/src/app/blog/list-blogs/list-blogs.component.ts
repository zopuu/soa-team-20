import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  ChangeDetectionStrategy,
} from '@angular/core';
import { BlogService } from '../blog.service';
import { Blog } from '../blog.model';
import { CommentService } from '../comment.service';
import { Like } from '../like.model';
import { LikeService } from '../like.service';
import { Comment } from '../comment.model';
import { AuthService } from 'src/app/auth/auth.service';

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
  // Map blogId to comments array
  blogComments: { [blogId: string]: Comment[] } = {};
  blogLikes: { [blogId: string]: Like[] } = {};
  // Map blogId to new comment text
  newCommentText: { [blogId: string]: string } = {};

  constructor(private blogService: BlogService,
    private commentService: CommentService,
    private likeService: LikeService,
    private authService: AuthService,
  ) {}

  ngOnInit() {
  this.authService.whoAmI().subscribe(user => {
    this.userId = user.id.toString(); // This is the logged-in user's ID
    // Use currentUserId as needed
  });
}
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
        // Fetch comments for each blog
        this.blogs.forEach(blog => {
          this.loadComments(blog.id);
          this.loadLikes(blog.id);});
      },
      error: (err) => {
        console.error('Failed to load blogs', err);
        this.error = 'Failed to load blogs.';
        this.loading = false;
      },
    });
  }
  loadComments(blogId: string) {
    console.log("Loading comments for blog:", blogId);
    this.commentService.getCommentsByBlog(blogId).subscribe({
      next: (comments) => {
        this.blogComments[blogId] = comments || [];
      },
      error: () => {
        this.blogComments[blogId] = [];
      }
    });
  }

  loadLikes(blogId: string) {
    console.log("Loading likes for blog:", blogId);
    this.likeService.getLikesByBlog(blogId).subscribe({
      next: (likes) => {
        this.blogLikes[blogId] = likes || [];
      },
      error: () => {
        this.blogLikes[blogId] = [];
      }
    });
  }

  hasLiked(blog: Blog): boolean {
  const likes = this.blogLikes[blog.id] || [];
  return likes.some(like => like.userId === this.userId);
}

toggleLike(blog: Blog) {
  if (this.hasLiked(blog) && this.userId) {
    // Unlike
    this.likeService.removeLike(blog.id, this.userId).subscribe({
      next: () => {
        this.blogLikes[blog.id] = (this.blogLikes[blog.id] || []).filter(like => like.userId !== this.userId);
      }
    });
  } else if(this.userId) {
    // Like
    this.likeService.createLike({ blogId: blog.id, userId: this.userId }).subscribe({
      next: (newLike) => {
        this.blogLikes[blog.id] = [...(this.blogLikes[blog.id] || []), newLike];
      }
    });
  }
}

  addComment(blog: Blog) {
    const text = this.newCommentText[blog.id];
    if (!text) return;
    const commentPayload = {
      blogId: blog.id,
      userId: this.userId,
      text,
      // userId: ... // set current user id here
    };
    console.log("Adding comment:", commentPayload);
    this.commentService.addComment(commentPayload).subscribe({
      next: () => {
        this.newCommentText[blog.id] = '';
        this.loadComments(blog.id); // Refresh comments
      }
    });
  }

  trackById = (_: number, b: Blog) => b.id;
}
