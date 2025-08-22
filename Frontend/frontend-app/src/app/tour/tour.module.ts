import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { CreateTourComponent } from './create-tour/create-tour.component';
import { CreateKeypointComponent } from './create-keypoint/create-keypoint.component';

@NgModule({
  declarations: [CreateTourComponent, CreateKeypointComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
  ],
})
export class TourModule {}
