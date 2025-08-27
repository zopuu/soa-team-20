import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { CreateTourComponent } from './create-tour/create-tour.component';
import { CreateKeypointComponent } from './create-keypoint/create-keypoint.component';
import { ListToursComponent } from './list-tours/list-tours.component';
import { ViewTourComponent } from './view-tour/view-tour.component';
import { KeypointDetailDialogComponent } from './keypoint-detail-dialog/keypoint-detail-dialog.component';
import { KeypointEditDialogComponent } from './keypoint-edit-dialog/keypoint-edit-dialog.component';
import { MapPickerDialogComponent } from './map-picker-dialog/map-picker-dialog.component';

@NgModule({
  declarations: [
    CreateTourComponent,
    CreateKeypointComponent,
    ListToursComponent,
    ViewTourComponent,
    // dialog used to show keypoint details on hover
    KeypointDetailDialogComponent,
    KeypointEditDialogComponent,
    MapPickerDialogComponent,
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatDialogModule,
    MatSelectModule,
    MatCheckboxModule,
  ],
})
export class TourModule {}
