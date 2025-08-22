import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-create-tour',
  templateUrl: './create-tour.component.html',
  styleUrls: ['./create-tour.component.css'],
})
export class CreateTourComponent {
  tourForm: FormGroup;

  constructor(private fb: FormBuilder) {
    this.tourForm = this.fb.group({
      name: ['', Validators.required],
      description: ['', [Validators.required, Validators.minLength(10)]],
      image: [''],
    });
  }

  onSubmit(): void {
    if (this.tourForm.valid) {
      const newTour = {
        name: this.tourForm.value.name,
        description: this.tourForm.value.description,
        image: this.tourForm.value.image,
      };
      // TODO: Call tour service to create tour
      console.log('Tour created:', newTour);
      this.tourForm.reset();
    }
  }
}
