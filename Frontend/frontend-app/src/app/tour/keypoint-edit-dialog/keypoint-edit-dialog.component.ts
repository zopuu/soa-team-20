import {
  Component,
  Inject,
  OnInit,
  ViewChild,
  ElementRef,
  AfterViewInit,
} from '@angular/core';
import {
  MAT_DIALOG_DATA,
  MatDialogRef,
  MatDialog,
} from '@angular/material/dialog';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { KeypointService } from '../keypoint.service';
import { MapPickerDialogComponent } from '../map-picker-dialog/map-picker-dialog.component';

@Component({
  selector: 'app-keypoint-edit-dialog',
  templateUrl: './keypoint-edit-dialog.component.html',
  styleUrls: ['./keypoint-edit-dialog.component.css'],
})
export class KeypointEditDialogComponent implements OnInit {
  form: FormGroup;
  selectedImage?: File;
  currentImageUrl?: string;

  constructor(
    private fb: FormBuilder,
    private dialogRef: MatDialogRef<KeypointEditDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private keypointService: KeypointService,
    private dialog: MatDialog
  ) {
    this.form = this.fb.group({
      title: ['', Validators.required],
      description: [''],
      latitude: [{ value: null, disabled: true }, Validators.required],
      longitude: [{ value: null, disabled: true }, Validators.required],
    });
  }

  ngOnInit(): void {
    if (this.data) {
      this.form.patchValue({
        title: this.data.title,
        description: this.data.description,
      });
      if (this.data.coordinates) {
        this.form.get('latitude')?.setValue(this.data.coordinates.latitude);
        this.form.get('longitude')?.setValue(this.data.coordinates.longitude);
      }
      // Set current image URL if keypoint has binary image data
      if (this.data.image?.data && this.data.image?.mimeType) {
        this.currentImageUrl = this.createImageDataUrl(
          this.data.image.data,
          this.data.image.mimeType
        );
      }
    }
  }

  private createImageDataUrl(base64Data: string, mimeType: string): string {
    // If the base64 data doesn't have the data URL prefix, add it
    if (base64Data.startsWith('data:')) {
      return base64Data;
    }
    return `data:${mimeType};base64,${base64Data}`;
  }

  onImageSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.selectedImage = input.files[0];
    }
  }

  onChangeCoordinates() {
    const current = {
      latitude: this.form.get('latitude')?.value,
      longitude: this.form.get('longitude')?.value,
    };
    const ref = this.dialog.open(MapPickerDialogComponent, {
      data: current,
      width: '90vw',
      height: '80vh',
      maxWidth: '95vw',
      maxHeight: '90vh',
      panelClass: 'map-picker-panel',
    });
    ref.afterClosed().subscribe((res) => {
      if (!res) return;
      if (res.latitude != null && res.longitude != null) {
        this.form.get('latitude')?.setValue(res.latitude);
        this.form.get('longitude')?.setValue(res.longitude);
      }
    });
  }

  onCancel() {
    this.dialogRef.close({ action: 'cancel' });
  }

  onSave() {
    if (!this.data?.id) return;
    const dto: any = {
      tourId: this.data.tourId, // Include tourId for backend
      title: this.form.get('title')?.value,
      description: this.form.get('description')?.value,
      coordinates: {
        latitude: Number(this.form.get('latitude')?.value),
        longitude: Number(this.form.get('longitude')?.value),
      },
    };
    this.keypointService
      .update(this.data.id, dto, this.selectedImage)
      .subscribe({
        next: (kp) => this.dialogRef.close({ action: 'updated', kp }),
        error: () => this.dialogRef.close({ action: 'update-failed' }),
      });
  }
}
