import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BlogService } from '../blog.service';
import { AuthService } from 'src/app/auth/auth.service';
import { BlogDto } from '../blog.dto';
import { firstValueFrom } from 'rxjs';

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

  async onSubmit(): Promise<void> {
  if (this.blogForm.valid) {
    try {
      const user = await firstValueFrom(this.auth.whoAmI());

      const newBlog: BlogDto = {
        userId: user.id.toString(),
        title: this.blogForm.value.title,
        description: this.blogForm.value.description,
      };

      await firstValueFrom(this.blogService.create(newBlog, this.selectedImages));

      console.log('Blog created successfully!');
      console.log('selected images:', this.selectedImages);
      this.blogForm.reset();
      this.selectedImages = [];
    } catch (err) {
      console.error('Error creating blog:', err);
    }
  }
}
}
