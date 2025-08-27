import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { KeypointService } from '../keypoint.service';
import { MatDialog } from '@angular/material/dialog';
import { KeypointEditDialogComponent } from '../keypoint-edit-dialog/keypoint-edit-dialog.component';

@Component({
  selector: 'app-keypoint-detail-dialog',
  template: `
    <h3 mat-dialog-title>{{ data?.title || 'Keypoint' }}</h3>
    <div mat-dialog-content>
      <p *ngIf="data?.description">{{ data.description }}</p>
      <p *ngIf="data?.coordinates">
        Lat: {{ data.coordinates.latitude }}, Lng:
        {{ data.coordinates.longitude }}
      </p>
    </div>
    <div mat-dialog-actions>
      <button mat-button (click)="onUpdate()">Update</button>
      <button mat-button color="warn" (click)="onDelete()">Delete</button>
      <button mat-button (click)="onClose()">Close</button>
    </div>
  `,
})
export class KeypointDetailDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<KeypointDetailDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private router: Router,
    private keypointService: KeypointService,
    private dialog: MatDialog
  ) {}

  onClose() {
    this.dialogRef.close({ action: 'close' });
  }

  onUpdate() {
    // open inline edit dialog and propagate result to parent
    const editRef = this.dialog.open(KeypointEditDialogComponent, {
      data: this.data,
      width: '640px',
    });
    editRef.afterClosed().subscribe((res) => {
      if (res && res.action === 'updated') {
        // close this detail dialog and notify parent that an update happened
        this.dialogRef.close({ action: 'updated', kp: res.kp });
      }
    });
  }

  onDelete() {
    if (!this.data?.id) {
      this.dialogRef.close({ action: 'delete' });
      return;
    }
    this.keypointService.delete(this.data.id).subscribe({
      next: () => this.dialogRef.close({ action: 'deleted', id: this.data.id }),
      error: () => this.dialogRef.close({ action: 'delete-failed' }),
    });
  }
}
