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
  selectedImages: File[] = [];

  constructor(
    private fb: FormBuilder,
    private blogService: BlogService,
    private auth: AuthService
  ) {
    this.blogForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', [Validators.required, Validators.minLength(10)]],
    });
  }

  onImagesSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files) {
      this.selectedImages = Array.from(input.files);
    }
  }

  onSubmit(): void {
    if (this.blogForm.valid) {
      const newBlog: BlogDto = {
        userId: '1',
        title: this.blogForm.value.title,
        description: this.blogForm.value.description,
      };

      this.blogService.create(newBlog, this.selectedImages).subscribe({
        next: () => {
          console.log('Blog created successfully!');
          console.log('selected images:', this.selectedImages);
          this.blogForm.reset();
          this.selectedImages = [];
        },
        error: (err) => {
          console.error('Error creating blog:', err);
        },
      });
    }
  }
}
