import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { KeypointService } from '../keypoint.service';
import { MatDialog } from '@angular/material/dialog';
import { KeypointEditDialogComponent } from '../keypoint-edit-dialog/keypoint-edit-dialog.component';

@Component({
  selector: 'app-keypoint-detail-dialog',
  templateUrl: './keypoint-detail-dialog.component.html',
  styleUrls: ['./keypoint-detail-dialog.component.css'],
})
export class KeypointDetailDialogComponent {
  imageUrl?: string;
  imageError = false;

  constructor(
    public dialogRef: MatDialogRef<KeypointDetailDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private router: Router,
    private keypointService: KeypointService,
    private dialog: MatDialog
  ) {
    // Set image URL if keypoint has binary image data
    if (data?.image?.data && data?.image?.mimeType) {
      this.imageUrl = this.createImageDataUrl(
        data.image.data,
        data.image.mimeType
      );
      console.log('Image data found:', data.image);
    }
  }

  private createImageDataUrl(base64Data: string, mimeType: string): string {
    // If the base64 data doesn't have the data URL prefix, add it
    if (base64Data.startsWith('data:')) {
      return base64Data;
    }
    return `data:${mimeType};base64,${base64Data}`;
  }

  onImageError(): void {
    this.imageError = true;
    console.error('Failed to load image');
  }

  onImageLoad(): void {
    this.imageError = false;
  }

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
