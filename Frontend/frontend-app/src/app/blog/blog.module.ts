import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CreateBlogComponent } from './create-blog/create-blog.component';
import { ListBlogsComponent } from './list-blogs/list-blogs.component';
import { MatFormFieldModule } from '@angular/material/form-field';
import { ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';

@NgModule({
  declarations: [CreateBlogComponent, ListBlogsComponent],
  imports: [
    CommonModule,
    MatFormFieldModule,
    ReactiveFormsModule,
    MatInputModule,
    MatButtonModule,
  ],
  exports: [ListBlogsComponent],
})
export class BlogModule {}
