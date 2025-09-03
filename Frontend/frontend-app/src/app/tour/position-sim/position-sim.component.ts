import { AfterViewInit, Component, ElementRef, OnDestroy, ViewChild } from '@angular/core';
import * as L from 'leaflet';
import { TouristLocationService, Coordinates } from '../tourist-location.service';
import { AuthService } from '../../auth/auth.service';

@Component({
  selector: 'app-position-sim',
  templateUrl: './position-sim.component.html',
  styleUrls: ['./position-sim.component.css']
})
export class PositionSimComponent implements AfterViewInit, OnDestroy {
  @ViewChild('map', { static: true }) mapEl!: ElementRef<HTMLDivElement>;

  private map?: L.Map;
  private marker?: L.Marker;
  private icon = L.icon({
    iconUrl: 'assets/icons/positionmark.svg',
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  userId?: number;
  selected?: Coordinates;
  statusMsg = '';      // shows inline “Saving… / Saved / Failed”
  hasKnownPosition = false; // for the initial hint

  constructor(
    private api: TouristLocationService,
    private auth: AuthService
  ) {}

  ngAfterViewInit(): void {
    this.auth.whoAmI().subscribe({
      next: (me) => {
        this.userId = me.id;
        this.initMap();
        this.loadExisting();
      },
      error: () => {
        this.statusMsg = 'Cannot determine user identity.';
      }
    });
  }

  private initMap() {
    this.map = L.map(this.mapEl.nativeElement, { center: [44.7866, 20.4489], zoom: 6 });
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // Click → set marker → auto-save
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      const coords: Coordinates = {
        latitude: Number(e.latlng.lat.toFixed(6)),
        longitude: Number(e.latlng.lng.toFixed(6))
      };
      this.placeMarker(coords);
      this.autoSave(coords);
    });
  }

  private loadExisting() {
    if (!this.userId) return;
    this.api.get(this.userId).subscribe({
      next: (loc) => {
        if (loc && this.map) {
          const { latitude, longitude } = loc.coordinates;
          this.placeMarker({ latitude, longitude });
          this.map.setView([latitude, longitude], 13);
          this.hasKnownPosition = true;
        } else {
          this.hasKnownPosition = false; // show the hint
        }
      },
      error: () => { this.hasKnownPosition = false; }
    });
  }

  private placeMarker(coords: Coordinates) {
    if (!this.map) return;
    if (this.marker) this.marker.setLatLng([coords.latitude, coords.longitude]);
    else this.marker = L.marker([coords.latitude, coords.longitude], { icon: this.icon }).addTo(this.map);
    this.selected = coords;
  }

  private autoSave(coords: Coordinates) {
    if (!this.userId) return;
    this.statusMsg = 'Saving…';
    this.api.set(this.userId, coords).subscribe({
      next: () => {
        this.statusMsg = 'Saved.';
        this.hasKnownPosition = true;
        // Clear the status after a short delay (optional)
        setTimeout(() => (this.statusMsg = ''), 1500);
      },
      error: () => {
        this.statusMsg = 'Failed to save location.';
      }
    });
  }

  ngOnDestroy(): void {
    if (this.map) { this.map.remove(); this.map = undefined; }
  }
}
