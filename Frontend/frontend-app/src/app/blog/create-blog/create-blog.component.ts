import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BlogService } from '../blog.service';
import { AuthService } from 'src/app/auth/auth.service';
import { BlogDto } from '../blog.dto';

@Component({
  selector: 'app-create-blog',
  templateUrl: './create-blog.component.html',
  styleUrls: ['./create-blog.component.css'],
})
export class CreateBlogComponent {
  blogForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private blogService: BlogService,
    private auth: AuthService
  ) {
    this.blogForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', [Validators.required, Validators.minLength(10)]],
      images: [''],
    });
  }

  onSubmit(): void {
    if (this.blogForm.valid) {
      const newBlog: BlogDto = {
        userId: '1',
        title: this.blogForm.value.title,
        description: this.blogForm.value.description,
        images: this.blogForm.value.images
          ? this.blogForm.value.images.split(',')
          : [],
      };

      this.blogService.create(newBlog).subscribe({
        next: () => {
          console.log('Blog created successfully!');
          this.blogForm.reset();
        },
        error: (err) => {
          console.error('Error creating blog:', err);
        },
      });
    }
  }
}
