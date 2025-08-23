import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { CreateTourComponent } from './create-tour/create-tour.component';
import { CreateKeypointComponent } from './create-keypoint/create-keypoint.component';
import { ListToursComponent } from './list-tours/list-tours.component';

@NgModule({
  declarations: [
    CreateTourComponent,
    CreateKeypointComponent,
    ListToursComponent,
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatCheckboxModule,
  ],
})
export class TourModule {}
