import {
  Component,
  Inject,
  AfterViewInit,
  ViewChild,
  ElementRef,
  OnDestroy,
} from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import * as L from 'leaflet';

@Component({
  selector: 'app-map-picker-dialog',
  template: `
    <h3 mat-dialog-title>Pick coordinates</h3>
    <div mat-dialog-content class="map-picker-root">
      <div #mapContainer class="map-container"></div>
      <div class="selected-row" *ngIf="selected">
        Selected: {{ selected.latitude }}, {{ selected.longitude }}
      </div>
    </div>
    <div mat-dialog-actions>
      <button mat-button (click)="onCancel()">Cancel</button>
      <button
        mat-button
        color="primary"
        (click)="onSave()"
        [disabled]="!selected"
      >
        Save
      </button>
    </div>
  `,
  styles: [
    `
      .map-picker-root {
        display: flex;
        flex-direction: column;
        padding: 0;
        margin: 0;
        height: calc(80vh - 96px);
      }
      .map-container {
        flex: 1 1 auto;
        width: 100%;
        height: 100%;
      }
      .selected-row {
        padding: 8px 0 12px 0;
      }
      :host ::ng-deep .mat-dialog-content {
        padding: 0 24px 12px 24px;
        overflow: hidden;
      }
      :host ::ng-deep .mat-dialog-container.map-picker-panel {
        padding: 0;
      }
    `,
  ],
})
export class MapPickerDialogComponent implements AfterViewInit, OnDestroy {
  @ViewChild('mapContainer', { static: true })
  mapContainer!: ElementRef<HTMLDivElement>;
  private map?: L.Map;
  private marker?: L.Marker;
  selected: { latitude: number; longitude: number } | null = null;

  constructor(
    public dialogRef: MatDialogRef<MapPickerDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}

  ngAfterViewInit(): void {
    const center: [number, number] = [
      this.data?.latitude ?? 20,
      this.data?.longitude ?? 0,
    ];
    this.map = L.map(this.mapContainer.nativeElement, {
      center,
      zoom: this.data?.latitude ? 12 : 2,
    });
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(this.map);

    // use project icon for the picker marker (position mark)
    const pickerIcon = L.icon({
      iconUrl: 'assets/icons/positionmark.svg',
      iconSize: [32, 32],
      iconAnchor: [16, 32],
    });

    if (this.data?.latitude && this.data?.longitude) {
      this.marker = L.marker([this.data.latitude, this.data.longitude], {
        icon: pickerIcon,
      }).addTo(this.map);
      this.selected = {
        latitude: this.data.latitude,
        longitude: this.data.longitude,
      };
    }

    this.map.on('click', (e: L.LeafletMouseEvent) => {
      const lat = e.latlng.lat;
      const lng = e.latlng.lng;
      this.setMarker(lat, lng, pickerIcon);
    });
  }

  private setMarker(lat: number, lng: number, icon?: L.Icon) {
    if (!this.map) return;
    if (this.marker) {
      this.marker.setLatLng([lat, lng]);
      if (icon) this.marker.setIcon(icon);
    } else {
      this.marker = L.marker([lat, lng], icon ? { icon } : undefined).addTo(
        this.map
      );
    }
    this.selected = {
      latitude: Number(lat.toFixed(6)),
      longitude: Number(lng.toFixed(6)),
    };
  }

  onCancel() {
    this.dialogRef.close();
  }

  onSave() {
    this.dialogRef.close(this.selected);
  }

  ngOnDestroy(): void {
    if (this.map) {
      this.map.remove();
      this.map = undefined;
    }
  }
}
