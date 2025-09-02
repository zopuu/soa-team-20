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
  template: `
    <h3 mat-dialog-title>Edit keypoint</h3>
    <form [formGroup]="form">
      <div mat-dialog-content>
        <mat-form-field appearance="fill" style="width:100%">
          <mat-label>Title</mat-label>
          <input matInput formControlName="title" />
        </mat-form-field>

        <mat-form-field appearance="fill" style="width:100%">
          <mat-label>Description</mat-label>
          <textarea matInput formControlName="description"></textarea>
        </mat-form-field>

        <mat-form-field appearance="fill" style="width:100%">
          <mat-label>Image URL</mat-label>
          <input matInput formControlName="image" />
        </mat-form-field>

        <div style="display:flex;gap:8px;align-items:center;margin-top:8px">
          <mat-form-field appearance="fill" style="flex:1">
            <mat-label>Latitude</mat-label>
            <input matInput formControlName="latitude" readonly />
          </mat-form-field>
          <mat-form-field appearance="fill" style="flex:1">
            <mat-label>Longitude</mat-label>
            <input matInput formControlName="longitude" readonly />
          </mat-form-field>
          <button mat-button type="button" (click)="onChangeCoordinates()">
            Change coordinates
          </button>
        </div>
      </div>
      <div mat-dialog-actions>
        <button mat-button (click)="onCancel()">Cancel</button>
        <button
          mat-button
          color="primary"
          (click)="onSave()"
          [disabled]="form.invalid"
        >
          Save
        </button>
      </div>
    </form>
  `,
})
export class KeypointEditDialogComponent implements OnInit {
  form: FormGroup;

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
      image: [''],
      latitude: [{ value: null, disabled: true }, Validators.required],
      longitude: [{ value: null, disabled: true }, Validators.required],
    });
  }

  ngOnInit(): void {
    if (this.data) {
      this.form.patchValue({
        title: this.data.title,
        description: this.data.description,
        image: this.data.image,
      });
      if (this.data.coordinates) {
        this.form.get('latitude')?.setValue(this.data.coordinates.latitude);
        this.form.get('longitude')?.setValue(this.data.coordinates.longitude);
      }
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
      title: this.form.get('title')?.value,
      description: this.form.get('description')?.value,
      image: this.form.get('image')?.value,
      coordinates: {
        latitude: Number(this.form.get('latitude')?.value),
        longitude: Number(this.form.get('longitude')?.value),
      },
    };
    this.keypointService.update(this.data.id, dto).subscribe({
      next: (kp) => this.dialogRef.close({ action: 'updated', kp }),
      error: () => this.dialogRef.close({ action: 'update-failed' }),
    });
  }
}
